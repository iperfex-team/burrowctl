package client

// RPCResponse represents the response structure received from the burrowctl server.
// This structure defines the standardized format for all server responses,
// regardless of the operation type (SQL queries, function calls, or system commands).
//
// The response follows a consistent tabular format where:
// - All results are represented as tables with named columns
// - Each row contains values corresponding to the column definitions
// - Errors are reported in a dedicated error field
//
// This design enables uniform handling of diverse operation types while
// maintaining compatibility with Go's database/sql interface expectations.
type RPCResponse struct {
	Columns []string        `json:"columns"` // Column names for the result table
	Rows    [][]interface{} `json:"rows"`    // Data rows, each containing values for all columns
	Error   string          `json:"error"`   // Error message if operation failed (empty on success)
}
