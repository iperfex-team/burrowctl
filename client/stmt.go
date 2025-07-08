package client

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// Stmt implements the database/sql/driver.Stmt interface for prepared statements.
// It provides efficient execution of the same query with different parameters,
// reducing parsing overhead and improving security through parameter binding.
//
// The statement maintains:
// - Reference to the parent connection
// - Original query string for execution
// - Parameter count for validation
// - Prepared statement lifecycle management
type Stmt struct {
	conn     *Conn  // Parent connection for execution
	query    string // Original SQL query with placeholders
	numInput int    // Number of placeholder parameters in the query
	closed   bool   // Whether the statement has been closed
}

// Close implements the driver.Stmt interface and releases statement resources.
// After closing, the statement cannot be executed again.
//
// Returns:
//   - error: Always nil as no special cleanup is required for RabbitMQ statements
func (s *Stmt) Close() error {
	s.closed = true
	s.conn.logf("Prepared statement closed: %s", s.query)
	return nil
}

// NumInput implements the driver.Stmt interface and returns the number of
// placeholder parameters in the prepared statement.
//
// Returns:
//   - int: Number of parameter placeholders (?) in the original query
func (s *Stmt) NumInput() int {
	return s.numInput
}

// Exec implements the driver.Stmt interface for executing prepared statements
// that don't return rows (INSERT, UPDATE, DELETE, etc.).
//
// Parameters:
//   - args: Parameter values to bind to the query placeholders
//
// Returns:
//   - driver.Result: Execution result (typically showing affected rows)
//   - error: Any error that occurred during execution
func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	if s.closed {
		return nil, fmt.Errorf("statement is closed")
	}

	// Validate parameter count
	if len(args) != s.numInput {
		return nil, fmt.Errorf("expected %d parameters, got %d", s.numInput, len(args))
	}

	s.conn.logf("Executing prepared statement with %d parameters", len(args))

	// Convert to NamedValue for internal execution
	named := make([]driver.NamedValue, len(args))
	for i, v := range args {
		named[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), s.conn.config.Timeout)
	defer cancel()

	// Execute through existing RPC mechanism
	rows, err := s.conn.queryRPC(ctx, s.query, named)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// For non-SELECT statements, return a simple result
	// Note: In a real implementation, the server should return affected row counts
	return &Result{}, nil
}

// Query implements the driver.Stmt interface for executing prepared statements
// that return rows (SELECT statements).
//
// Parameters:
//   - args: Parameter values to bind to the query placeholders
//
// Returns:
//   - driver.Rows: Result set for iteration
//   - error: Any error that occurred during execution
func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	if s.closed {
		return nil, fmt.Errorf("statement is closed")
	}

	// Validate parameter count
	if len(args) != s.numInput {
		return nil, fmt.Errorf("expected %d parameters, got %d", s.numInput, len(args))
	}

	s.conn.logf("Querying prepared statement with %d parameters", len(args))

	// Convert to NamedValue for internal execution
	named := make([]driver.NamedValue, len(args))
	for i, v := range args {
		named[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), s.conn.config.Timeout)
	defer cancel()

	// Execute through existing RPC mechanism
	return s.conn.queryRPC(ctx, s.query, named)
}

// ExecContext implements the driver.StmtExecContext interface for context-aware
// execution of prepared statements that don't return rows.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - args: Parameter values to bind to the query placeholders
//
// Returns:
//   - driver.Result: Execution result (typically showing affected rows)
//   - error: Any error that occurred during execution
func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s.closed {
		return nil, fmt.Errorf("statement is closed")
	}

	// Validate parameter count
	if len(args) != s.numInput {
		return nil, fmt.Errorf("expected %d parameters, got %d", s.numInput, len(args))
	}

	s.conn.logf("Executing prepared statement (context) with %d parameters", len(args))

	// Execute through existing RPC mechanism
	rows, err := s.conn.queryRPC(ctx, s.query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// For non-SELECT statements, return a simple result
	return &Result{}, nil
}

// QueryContext implements the driver.StmtQueryContext interface for context-aware
// execution of prepared statements that return rows.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - args: Parameter values to bind to the query placeholders
//
// Returns:
//   - driver.Rows: Result set for iteration
//   - error: Any error that occurred during execution
func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s.closed {
		return nil, fmt.Errorf("statement is closed")
	}

	// Validate parameter count
	if len(args) != s.numInput {
		return nil, fmt.Errorf("expected %d parameters, got %d", s.numInput, len(args))
	}

	s.conn.logf("Querying prepared statement (context) with %d parameters", len(args))

	// Execute through existing RPC mechanism
	return s.conn.queryRPC(ctx, s.query, args)
}

// Result implements the driver.Result interface for prepared statement execution results.
// This is a simple implementation that would be enhanced with actual server-side
// result metadata in a production system.
type Result struct {
	affectedRows int64
	lastInsertID int64
}

// LastInsertId implements the driver.Result interface.
// Returns the last insert ID for INSERT statements.
//
// Returns:
//   - int64: Last insert ID (currently always 0)
//   - error: Always nil in this implementation
func (r *Result) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

// RowsAffected implements the driver.Result interface.
// Returns the number of rows affected by the statement.
//
// Returns:
//   - int64: Number of affected rows (currently always 0)
//   - error: Always nil in this implementation
func (r *Result) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}

// countPlaceholders counts the number of parameter placeholders (?) in a SQL query.
// This is used to validate the correct number of parameters are provided.
//
// Parameters:
//   - query: SQL query string to analyze
//
// Returns:
//   - int: Number of ? placeholders found in the query
func countPlaceholders(query string) int {
	count := 0
	inString := false
	escaped := false

	for _, char := range query {
		switch {
		case escaped:
			escaped = false
		case char == '\\':
			escaped = true
		case char == '\'' && !escaped:
			inString = !inString
		case char == '?' && !inString && !escaped:
			count++
		}
	}

	return count
}