package client

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
)

// Rows implements the database/sql/driver.Rows interface for burrowctl query results.
// It provides iteration over result sets returned from SQL queries, function calls,
// or system commands executed through the RabbitMQ transport layer.
//
// The Rows structure maintains:
// - Column names for the result set
// - Row data as a slice of interface{} slices
// - Current position for iteration
// - Type conversion logic for database values
//
// This implementation handles the conversion between JSON-serialized data
// from the server and Go's database/sql driver value types.
type Rows struct {
	columns []string        // Column names from the query result
	rows    [][]interface{} // Row data as received from server
	pos     int             // Current position in the result set
}

// Columns implements the driver.Rows interface and returns the column names
// for the current result set. This method is called by the database/sql
// package to populate sql.Rows.Columns().
//
// Returns:
//   - []string: Array of column names in order
func (r *Rows) Columns() []string {
	return r.columns
}

// Next implements the driver.Rows interface and advances to the next row.
// It populates the provided destination slice with converted values from
// the current row and advances the internal position counter.
//
// This method handles type conversion from JSON-serialized values to
// appropriate Go types for the database/sql interface.
//
// Parameters:
//   - dest: Destination slice to populate with row values
//
// Returns:
//   - error: io.EOF when no more rows are available, or any conversion error
func (r *Rows) Next(dest []driver.Value) error {
	// Check if we've reached the end of the result set
	if r.pos >= len(r.rows) {
		return io.EOF
	}

	// Convert and copy current row values to destination
	for i, val := range r.rows[r.pos] {
		dest[i] = r.convertValue(val)
	}

	// Advance to next row
	r.pos++
	return nil
}

// convertValue converts server response values to appropriate Go driver types.
// This method handles the complexity of converting JSON-deserialized values
// (which have limited type information) to the rich type system expected
// by Go's database/sql drivers.
//
// The conversion strategy:
// - Attempts to parse string representations of numbers back to numeric types
// - Converts JSON float64 values to int64 when they represent whole numbers
// - Preserves boolean values as-is
// - Converts unknown types to string representations
//
// Parameters:
//   - val: Raw value from server response (JSON-deserialized)
//
// Returns:
//   - driver.Value: Converted value suitable for database/sql interface
func (r *Rows) convertValue(val interface{}) driver.Value {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		// Attempt to convert string representations of numbers back to numeric types
		// This handles cases where the server sends numbers as strings for precision
		if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			return intVal
		}
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal
		}
		// Return as string if not a number
		return v
	case float64:
		// JSON unmarshaling always returns float64 for numbers
		// Convert to int64 if it represents a whole number
		if v == float64(int64(v)) {
			return int64(v)
		}
		return v
	case bool:
		// Boolean values pass through unchanged
		return v
	default:
		// Convert unknown types to string representation
		return fmt.Sprintf("%v", v)
	}
}

// Close implements the driver.Rows interface and cleans up any resources.
// For the burrowctl client, no special cleanup is required as all data
// is already in memory from the RPC response.
//
// Returns:
//   - error: Always nil as no cleanup is required
func (r *Rows) Close() error {
	return nil
}
