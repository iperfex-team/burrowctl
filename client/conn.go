package client

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net"
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
type Conn struct {
	deviceID  string             // Target device/server identifier
	connMgr   *ConnectionManager // Connection manager with automatic reconnection
	config    *DSNConfig         // Parsed DSN configuration
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
	return c.connMgr.Close()
}

// Begin implements the driver.Conn interface but transactions are not supported.
// The burrowctl system operates on individual message-based requests rather
// than transaction-based operations due to its distributed nature.
//
// Returns:
//   - driver.Tx: Always nil
//   - error: Always returns "transactions not supported" error
func (c *Conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("transactions not supported")
}

// Query implements the driver.Conn interface for executing queries with arguments.
// This method converts driver.Value arguments to driver.NamedValue and delegates
// to QueryContext with a timeout context.
//
// Parameters:
//   - query: SQL query, function call, or system command
//   - args: Query arguments as driver.Value slice
//
// Returns:
//   - driver.Rows: Result set from the query execution
//   - error: Any error that occurred during query execution
func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	startTotal := time.Now()
	c.logf("Executing query: %s", query)

	// Create timeout context based on configuration
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// Convert driver.Value arguments to driver.NamedValue format
	named := make([]driver.NamedValue, len(args))
	for i, v := range args {
		named[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}

	// Execute query through RPC
	rows, err := c.queryRPC(ctx, query, named)

	// Log total execution time if debug enabled
	total := time.Since(startTotal)
	c.logf("total time: %v", total)

	return rows, err
}

// QueryContext implements the driver.ConnContext interface for context-aware query execution.
// This is the primary method for executing all types of operations (SQL, functions, commands)
// with proper timeout and cancellation support.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - query: SQL query, function call, or system command
//   - args: Query arguments as driver.NamedValue slice
//
// Returns:
//   - driver.Rows: Result set from the query execution
//   - error: Any error that occurred during query execution
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	startTotal := time.Now()
	c.logf("Executing query (Context): %s", query)

	// Execute query through RPC
	rows, err := c.queryRPC(ctx, query, args)

	// Log total execution time if debug enabled
	total := time.Since(startTotal)
	c.logf("total time (QueryContext): %v", total)

	return rows, err
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

// queryRPC executes a query through the RabbitMQ RPC mechanism.
// This is the core method that handles the complete request-response cycle
// including message routing, timeout handling, and response processing.
//
// The method implements the RabbitMQ RPC pattern:
// 1. Creates a temporary reply queue for receiving responses
// 2. Publishes the query to the device-specific queue
// 3. Waits for the response with timeout handling
// 4. Processes and validates the response
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - query: The query string (SQL, function call, or command)
//   - args: Query arguments as driver.NamedValue slice
//
// Returns:
//   - driver.Rows: Result set containing query results
//   - error: Any error that occurred during RPC execution
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
		"type":     cmdType,                // Query type: sql, function, or command
		"deviceID": c.deviceID,             // Target device identifier
		"query":    actualQuery,            // Actual query without prefix
		"params":   argsToSlice(args),      // Query parameters
		"clientIP": getOutboundIP(),        // Client IP for logging
	}

	// Serialize request to JSON
	body, _ := json.Marshal(req)

	startRT := time.Now()
	c.logf("Publishing query to device queue '%s'", c.deviceID)

	// Publish query to device-specific queue with RPC headers
	err = ch.PublishWithContext(ctx, "", c.deviceID, false, false, amqp.Publishing{
		ContentType:   "application/json", // JSON content type
		CorrelationId: corrID,             // For matching request/response
		ReplyTo:       replyQueue.Name,    // Where to send the response
		Body:          body,               // Serialized request
	})
	if err != nil {
		return nil, fmt.Errorf("failed to publish query to device queue '%s': %v\nPlease check:\n- Server is running\n- Device ID '%s' is correct\n- Queue exists", c.deviceID, err, c.deviceID)
	}
	c.logf("Query published, waiting for response...")

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
