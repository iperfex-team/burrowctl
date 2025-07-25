package client

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Conn implements the database/sql/driver.Conn interface for RabbitMQ-based connections.
// It manages the connection state, handles query routing, and provides the core
// functionality for executing SQL queries, function calls, and system commands
// through the RabbitMQ transport layer.
//
// The connection maintains:
// - RabbitMQ connection for message transport with automatic reconnection
// - Device ID for routing messages to the correct server
// - Configuration including timeouts and debug settings
// - Query execution context and error handling
// - Heartbeat manager for connection monitoring
type Conn struct {
	deviceID       string             // Target device/server identifier
	connMgr        *ConnectionManager // Connection manager with automatic reconnection
	config         *DSNConfig         // Parsed DSN configuration
	currentTx      *Tx                // Current active transaction (if any)
	transactionMux sync.RWMutex       // Mutex for transaction state

	// Heartbeat management
	heartbeatManager *HeartbeatManager // Heartbeat manager for connection monitoring
	rpcActive        bool              // Whether RPC is currently active
	rpcMutex         sync.RWMutex      // Mutex for RPC state
}

// logf provides conditional debug logging based on the configuration.
// Debug messages are only logged when debug mode is enabled in the DSN configuration.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Arguments for the format string
func (c *Conn) logf(format string, args ...interface{}) {
	if c.config != nil && c.config.Debug {
		log.Printf("[client debug] "+format, args...)
	}
}

// Prepare implements the driver.Conn interface and creates a prepared statement.
// Prepared statements provide performance benefits and security through parameter binding.
// The statement can be executed multiple times with different parameters.
//
// Parameters:
//   - query: SQL query string with parameter placeholders (?)
//
// Returns:
//   - driver.Stmt: Prepared statement ready for execution
//   - error: Any error that occurred during statement preparation
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	c.logf("Preparing statement: %s", query)

	// Count parameter placeholders for validation
	numInput := countPlaceholders(query)

	stmt := &Stmt{
		conn:     c,
		query:    query,
		numInput: numInput,
		closed:   false,
	}

	c.logf("Statement prepared with %d parameters", numInput)
	return stmt, nil
}

// Close implements the driver.Conn interface and closes the RabbitMQ connection.
// This method should be called when the database connection is no longer needed
// to properly clean up resources.
//
// Returns:
//   - error: Any error that occurred during connection closure
func (c *Conn) Close() error {
	c.logf("Closing connection to RabbitMQ")

	// Stop heartbeat manager if active
	if c.heartbeatManager != nil {
		c.heartbeatManager.Stop()
	}

	return c.connMgr.Close()
}

// Begin implements the driver.Conn interface and starts a new transaction.
// The transaction provides basic ACID properties through server-side coordination.
//
// Returns:
//   - driver.Tx: New transaction instance
//   - error: Any error that occurred during transaction start
func (c *Conn) Begin() (driver.Tx, error) {
	c.transactionMux.Lock()
	defer c.transactionMux.Unlock()

	// Check if there's already an active transaction
	if c.currentTx != nil && c.currentTx.IsActive() {
		return nil, fmt.Errorf("transaction already in progress")
	}

	c.logf("Starting new transaction")

	// Create new transaction
	tx := newTransaction(c)

	// Send BEGIN command to server
	err := tx.executeTransactionCommand("BEGIN")
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	c.currentTx = tx
	c.logf("Transaction started: %s", tx.transactionID)
	return tx, nil
}

// Query implements the driver.Conn interface and executes a query with parameters.
// It supports parameter binding for security and type safety.
//
// Parameters:
//   - query: SQL query string with parameter placeholders (?)
//   - args: Query parameters to bind to placeholders
//
// Returns:
//   - driver.Rows: Result set from the query
//   - error: Any error that occurred during query execution
func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	c.logf("Executing query: %s with %d parameters", query, len(args))

	// Convert driver.Value to driver.NamedValue for consistency
	namedArgs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		namedArgs[i] = driver.NamedValue{Name: "", Ordinal: i + 1, Value: arg}
	}

	// Create timeout context based on configuration
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// Execute query through RPC with heartbeat activation
	return c.queryRPCWithHeartbeat(ctx, query, namedArgs)
}

// QueryContext implements the driver.QueryerContext interface and executes a query
// with parameters using a context for cancellation and timeout control.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - query: SQL query string with parameter placeholders (?)
//   - args: Query parameters to bind to placeholders
//
// Returns:
//   - driver.Rows: Result set from the query
//   - error: Any error that occurred during query execution
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.logf("Executing query with context: %s with %d parameters", query, len(args))

	// Execute query through RPC with heartbeat activation
	return c.queryRPCWithHeartbeat(ctx, query, args)
}

// getOutboundIP determines the client's outbound IP address by establishing
// a UDP connection to a public DNS server. This IP is included in requests
// for server-side logging and debugging purposes.
//
// The function uses Google's public DNS (8.8.8.8) as a destination to determine
// the local outbound interface without actually sending any data.
//
// Returns:
//   - string: The client's outbound IP address, or "unknown" if determination fails
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// parseCommand analyzes a query string to determine its type and extract the actual command.
// The burrowctl system supports three types of operations, each with a specific prefix:
//
// - FUNCTION: prefix indicates a function call with JSON parameters
// - COMMAND: prefix indicates a system command execution
// - No prefix: indicates a standard SQL query
//
// Parameters:
//   - query: The raw query string to analyze
//
// Returns:
//   - cmdType: The detected command type ("sql", "function", or "command")
//   - actualQuery: The query string with any prefix removed
//
// Examples:
//   - "SELECT * FROM users" → ("sql", "SELECT * FROM users")
//   - "FUNCTION:{"name":"test"}" → ("function", "{"name":"test"}")
//   - "COMMAND:ls -la" → ("command", "ls -la")
func parseCommand(query string) (cmdType string, actualQuery string) {
	// Check for function call prefix
	if len(query) > 9 && query[:9] == "FUNCTION:" {
		return "function", query[9:]
	}
	// Check for system command prefix
	if len(query) > 8 && query[:8] == "COMMAND:" {
		return "command", query[8:]
	}
	// Default to SQL query
	return "sql", query
}

// queryRPCWithHeartbeat executes RPC with heartbeat activation
func (c *Conn) queryRPCWithHeartbeat(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	// Activate heartbeat at the start of RPC
	c.activateHeartbeat()
	defer c.deactivateHeartbeat()

	// Execute the actual RPC
	return c.queryRPC(ctx, query, args)
}

// activateHeartbeat activates the heartbeat when RPC is active
func (c *Conn) activateHeartbeat() {
	c.rpcMutex.Lock()
	defer c.rpcMutex.Unlock()

	if !c.rpcActive {
		c.rpcActive = true
		if c.heartbeatManager != nil {
			c.heartbeatManager.ActivateHeartbeat()
		}
		c.logf("RPC activated, heartbeat enabled")
	}
}

// deactivateHeartbeat deactivates the heartbeat when RPC is not active
func (c *Conn) deactivateHeartbeat() {
	c.rpcMutex.Lock()
	defer c.rpcMutex.Unlock()

	if c.rpcActive {
		c.rpcActive = false
		if c.heartbeatManager != nil {
			c.heartbeatManager.DeactivateHeartbeat()
		}
		c.logf("RPC deactivated, heartbeat disabled")
	}
}

// queryRPC sends a query to the server via RabbitMQ RPC using separate RPC queue
func (c *Conn) queryRPC(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	// Get current connection from connection manager
	conn, err := c.connMgr.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("no active connection: %v", err)
	}

	// Create RabbitMQ channel for this query
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %v", err)
	}
	defer ch.Close()
	c.logf("RabbitMQ channel opened")

	// Declare exclusive reply queue for receiving response
	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %v", err)
	}
	c.logf("Reply queue declared: %s", replyQueue.Name)

	// Generate unique correlation ID for request-response matching
	corrID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Parse query to determine type and extract actual command
	cmdType, actualQuery := parseCommand(query)
	c.logf("Detected command type: %s, actual query: %s", cmdType, actualQuery)

	// Build RPC request message
	req := map[string]interface{}{
		"type":     cmdType,           // Query type: sql, function, or command
		"deviceID": c.deviceID,        // Target device identifier
		"query":    actualQuery,       // Actual query without prefix
		"params":   argsToSlice(args), // Query parameters
		"clientIP": getOutboundIP(),   // Client IP for logging
	}

	// Include transaction information if we're in a transaction
	c.transactionMux.RLock()
	if c.currentTx != nil && c.currentTx.IsActive() {
		req["transactionID"] = c.currentTx.GetTransactionID()
		c.logf("Query executing in transaction: %s", c.currentTx.GetTransactionID())
	}
	c.transactionMux.RUnlock()

	// Serialize request to JSON
	body, _ := json.Marshal(req)

	startRT := time.Now()
	c.logf("Publishing query to device RPC queue '%s'", c.deviceID)

	// Publish query to device-specific RPC queue (separate from heartbeat)
	rpcQueueName := fmt.Sprintf("device_%s_rpc", c.deviceID)
	err = ch.PublishWithContext(ctx, "", rpcQueueName, false, false, amqp.Publishing{
		ContentType:   "application/json", // JSON content type
		CorrelationId: corrID,             // For matching request/response
		ReplyTo:       replyQueue.Name,    // Where to send the response
		Body:          body,               // Serialized request
	})
	if err != nil {
		return nil, fmt.Errorf("failed to publish query to device RPC queue '%s': %v\nPlease check:\n- Server is running\n- Device ID '%s' is correct\n- Queue exists", rpcQueueName, err, c.deviceID)
	}
	c.logf("Query published to RPC queue, waiting for response...")

	// Start consuming from reply queue
	msgs, err := ch.Consume(replyQueue.Name, "", true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume from reply queue: %v", err)
	}

	// Wait for response or timeout
	select {
	case <-ctx.Done():
		// Context cancelled or timed out
		return nil, fmt.Errorf("timeout (%v) waiting for device response from '%s'\nPlease check:\n- Server is running and responding\n- Device ID '%s' is correct\n- Database is accessible", c.config.Timeout, c.deviceID, c.deviceID)
	case msg := <-msgs:
		// Response received
		rt := time.Since(startRT)
		c.logf("RabbitMQ roundtrip time: %v", rt)

		// Validate correlation ID to ensure response matches request
		if msg.CorrelationId != corrID {
			return nil, fmt.Errorf("correlation id mismatch: expected %s, got %s", corrID, msg.CorrelationId)
		}

		// Parse server response
		var resp RPCResponse
		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse server response: %v", err)
		}

		// Check for server-side errors
		if resp.Error != "" {
			return nil, fmt.Errorf("server error: %s", resp.Error)
		}

		// Return successful result set
		c.logf("Response received with %d rows", len(resp.Rows))
		return &Rows{columns: resp.Columns, rows: resp.Rows}, nil
	}
}

// argsToSlice converts driver.NamedValue arguments to a plain interface{} slice.
// This conversion is necessary for JSON marshaling of query parameters.
//
// Parameters:
//   - args: Array of named values from the driver interface
//
// Returns:
//   - []interface{}: Plain slice containing just the values for JSON serialization
func argsToSlice(args []driver.NamedValue) []interface{} {
	var out []interface{}
	for _, a := range args {
		out = append(out, a.Value)
	}
	return out
}

// clearFinishedTransaction clears the current transaction reference if it's no longer active.
// This method should be called after transaction completion to clean up resources.
func (c *Conn) clearFinishedTransaction() {
	c.transactionMux.Lock()
	defer c.transactionMux.Unlock()

	if c.currentTx != nil && !c.currentTx.IsActive() {
		c.logf("Clearing finished transaction: %s", c.currentTx.GetTransactionID())
		c.currentTx = nil
	}
}

// setupHeartbeat initializes the heartbeat manager
func (c *Conn) setupHeartbeat() {
	if c.config.HeartbeatEnabled {
		c.heartbeatManager = NewHeartbeatManager(
			c.connMgr,
			c.deviceID,
			getOutboundIP(),
			c.config.HeartbeatConfig,
		)
		c.heartbeatManager.SetCallbacks(c.handleDisconnect, c.handleReconnect)
	}
}

// handleDisconnect callback for heartbeat manager
func (c *Conn) handleDisconnect(err error) {
	c.logf("Connection considered dead: %v", err)
	// Trigger reconnection
	c.connMgr.Reconnect()
}

// handleReconnect callback for heartbeat manager
func (c *Conn) handleReconnect() {
	c.logf("Connection restored")
}

// GetHeartbeatStats returns heartbeat statistics
func (c *Conn) GetHeartbeatStats() HeartbeatStats {
	if c.heartbeatManager != nil {
		return c.heartbeatManager.GetStats()
	}
	return HeartbeatStats{}
}
