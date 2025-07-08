package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Tx implements the database/sql/driver.Tx interface for transaction support.
// It provides basic transaction functionality over RabbitMQ by maintaining
// transaction state and coordinating BEGIN/COMMIT/ROLLBACK operations.
//
// The transaction implementation:
// - Maintains a unique transaction ID for server-side coordination
// - Tracks transaction state (active, committed, rolled back)
// - Provides timeout handling for transaction operations
// - Supports nested transaction detection and prevention
type Tx struct {
	conn           *Conn           // Parent connection
	transactionID  string          // Unique transaction identifier
	state          TxState         // Current transaction state
	startTime      time.Time       // When transaction began
	mutex          sync.RWMutex    // Thread-safe state access
	ctx            context.Context // Context for cancellation
	cancel         context.CancelFunc
}

// TxState represents the current state of a transaction
type TxState int

const (
	TxActive TxState = iota
	TxCommitted
	TxRolledBack
)

// String returns a string representation of the transaction state
func (ts TxState) String() string {
	switch ts {
	case TxActive:
		return "active"
	case TxCommitted:
		return "committed"
	case TxRolledBack:
		return "rolled_back"
	default:
		return "unknown"
	}
}

// newTransaction creates a new transaction instance.
// It generates a unique transaction ID and initializes the transaction state.
//
// Parameters:
//   - conn: Parent connection for the transaction
//
// Returns:
//   - *Tx: New transaction instance ready for use
func newTransaction(conn *Conn) *Tx {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 5 minute transaction timeout
	
	tx := &Tx{
		conn:          conn,
		transactionID: fmt.Sprintf("tx_%d_%d", time.Now().Unix(), time.Now().UnixNano()),
		state:         TxActive,
		startTime:     time.Now(),
		ctx:           ctx,
		cancel:        cancel,
	}

	conn.logf("Transaction created: %s", tx.transactionID)
	return tx
}

// Commit implements the driver.Tx interface and commits the transaction.
// It sends a COMMIT command to the server and marks the transaction as committed.
//
// Returns:
//   - error: Any error that occurred during commit
func (tx *Tx) Commit() error {
	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if tx.state != TxActive {
		return fmt.Errorf("transaction is not active (state: %s)", tx.state)
	}

	tx.conn.logf("Committing transaction: %s", tx.transactionID)

	// Send COMMIT command to server
	err := tx.executeTransactionCommand("COMMIT")
	if err != nil {
		tx.conn.logf("Transaction commit failed: %s, error: %v", tx.transactionID, err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	tx.state = TxCommitted
	tx.cancel() // Cancel context to free resources
	
	duration := time.Since(tx.startTime)
	tx.conn.logf("Transaction committed successfully: %s (duration: %v)", tx.transactionID, duration)
	
	// Clear transaction reference from connection
	tx.conn.clearFinishedTransaction()
	return nil
}

// Rollback implements the driver.Tx interface and rolls back the transaction.
// It sends a ROLLBACK command to the server and marks the transaction as rolled back.
//
// Returns:
//   - error: Any error that occurred during rollback
func (tx *Tx) Rollback() error {
	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if tx.state != TxActive {
		return fmt.Errorf("transaction is not active (state: %s)", tx.state)
	}

	tx.conn.logf("Rolling back transaction: %s", tx.transactionID)

	// Send ROLLBACK command to server
	err := tx.executeTransactionCommand("ROLLBACK")
	if err != nil {
		tx.conn.logf("Transaction rollback failed: %s, error: %v", tx.transactionID, err)
		return fmt.Errorf("failed to rollback transaction: %v", err)
	}

	tx.state = TxRolledBack
	tx.cancel() // Cancel context to free resources
	
	duration := time.Since(tx.startTime)
	tx.conn.logf("Transaction rolled back successfully: %s (duration: %v)", tx.transactionID, duration)
	
	// Clear transaction reference from connection
	tx.conn.clearFinishedTransaction()
	return nil
}

// executeTransactionCommand sends a transaction command (BEGIN, COMMIT, ROLLBACK) to the server.
// This method handles the RabbitMQ communication for transaction control.
//
// Parameters:
//   - command: Transaction command to execute ("BEGIN", "COMMIT", "ROLLBACK")
//
// Returns:
//   - error: Any error that occurred during command execution
func (tx *Tx) executeTransactionCommand(command string) error {
	// Get current connection from connection manager
	conn, err := tx.conn.connMgr.GetConnection()
	if err != nil {
		return fmt.Errorf("no active connection: %v", err)
	}

	// Create RabbitMQ channel for this transaction command
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare exclusive reply queue for receiving response
	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare reply queue: %v", err)
	}

	// Generate unique correlation ID for request-response matching
	corrID := fmt.Sprintf("tx_%d", time.Now().UnixNano())

	// Build transaction command request
	req := map[string]interface{}{
		"type":          "transaction",           // Special type for transaction commands
		"deviceID":      tx.conn.deviceID,        // Target device identifier
		"transactionID": tx.transactionID,        // Transaction ID for server-side tracking
		"command":       command,                 // Transaction command (BEGIN, COMMIT, ROLLBACK)
		"clientIP":      getOutboundIP(),         // Client IP for logging
	}

	// Serialize request to JSON
	body, _ := json.Marshal(req)

	tx.conn.logf("Sending transaction command '%s' for transaction %s", command, tx.transactionID)

	// Publish command to device-specific queue with RPC headers
	err = ch.PublishWithContext(tx.ctx, "", tx.conn.deviceID, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		ReplyTo:       replyQueue.Name,
		Body:          body,
	})
	if err != nil {
		return fmt.Errorf("failed to publish transaction command: %v", err)
	}

	// Start consuming from reply queue
	msgs, err := ch.Consume(replyQueue.Name, "", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume from reply queue: %v", err)
	}

	// Create timeout context for transaction command
	cmdCtx, cancel := context.WithTimeout(tx.ctx, 30*time.Second)
	defer cancel()

	// Wait for response or timeout
	select {
	case <-cmdCtx.Done():
		return fmt.Errorf("timeout waiting for transaction command response")
	case msg := <-msgs:
		// Validate correlation ID
		if msg.CorrelationId != corrID {
			return fmt.Errorf("correlation id mismatch: expected %s, got %s", corrID, msg.CorrelationId)
		}

		// Parse server response
		var resp RPCResponse
		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			return fmt.Errorf("failed to parse server response: %v", err)
		}

		// Check for server-side errors
		if resp.Error != "" {
			return fmt.Errorf("server error: %s", resp.Error)
		}

		tx.conn.logf("Transaction command '%s' completed successfully for transaction %s", command, tx.transactionID)
		return nil
	}
}

// IsActive returns whether the transaction is still active
func (tx *Tx) IsActive() bool {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()
	return tx.state == TxActive
}

// GetState returns the current transaction state
func (tx *Tx) GetState() TxState {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()
	return tx.state
}

// GetTransactionID returns the unique transaction identifier
func (tx *Tx) GetTransactionID() string {
	return tx.transactionID
}

// GetDuration returns how long the transaction has been active
func (tx *Tx) GetDuration() time.Duration {
	return time.Since(tx.startTime)
}