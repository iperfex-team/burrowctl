# Extended Client Example

This example demonstrates the new **Extended BurrowClient** that provides a cleaner interface for different operation types while maintaining full compatibility with the standard `database/sql` interface.

## 🚀 Features

### Clean API Interface
- **`ExecCommand()`**: Execute system commands with structured results
- **`ExecFunction()`**: Call custom functions with typed parameters
- **Standard SQL methods**: Query, QueryRow, Exec, Begin, Prepare
- **Full compatibility**: Access to underlying `sql.DB` instance

### Structured Results
- **`CommandResult`**: Structured command execution results with stdout/stderr separation
- **`FunctionResult`**: Typed function results with execution metadata
- **Type safety**: Proper Go types for all operations

### Helper Functions
- **Parameter helpers**: `StringParam()`, `IntParam()`, `Float64Param()`, etc.
- **Error handling**: Comprehensive error information
- **Metadata**: Execution time, timestamps, and performance data

## 🔧 Usage

### Basic Connection
```go
// Create extended client
bc, err := client.NewBurrowClient(dsn)
if err != nil {
    log.Fatal("Failed to create client:", err)
}
defer bc.Close()
```

### System Commands
```go
// Execute system command
result, err := bc.ExecCommand("df -h")
if err != nil {
    log.Fatal("Command failed:", err)
}

fmt.Printf("Exit code: %d\n", result.ExitCode)
fmt.Printf("Stdout: %v\n", result.Stdout)
fmt.Printf("Stderr: %v\n", result.Stderr)
```

### Function Calls
```go
// Call custom function
result, err := bc.ExecFunction("addIntegers", 
    client.IntParam(10), 
    client.IntParam(20),
)
if err != nil {
    log.Fatal("Function failed:", err)
}

fmt.Printf("Result: %v\n", result.Result)
fmt.Printf("Duration: %s\n", result.Duration)
```

### SQL Operations
```go
// Standard SQL queries
rows, err := bc.Query("SELECT * FROM users WHERE active = ?", true)
if err != nil {
    log.Fatal("Query failed:", err)
}
defer rows.Close()

// Process results...
```

### Database/SQL Compatibility
```go
// Access underlying sql.DB instance
db := bc.DB()

// Use standard database/sql methods
rows, err := db.Query("SELECT * FROM table")
// ... standard usage
```

## 🏃 Running the Example

### Prerequisites
- RabbitMQ server running on localhost:5672
- burrowctl server running with device ID "extended-client-demo"

### Start the server
```bash
cd examples/server/basic
go run main.go
```

### Run the example
```bash
cd examples/client/extended-client-example
go run main.go
```

## 📊 Example Output

```
🔌 Testing connection...
✅ Connection successful!

📊 Executing SQL queries...
  • Simple SELECT query:
    Result: Hello, 42
  • Parameterized query:
    Result: Alice is 30 years old

🖥️ Executing system commands...
  • Executing: echo 'Hello from system command'
    Exit code: 0
    Executed at: 2024-01-15T10:30:00Z
    Stdout:
      Hello from system command

⚙️ Executing custom functions...
  • lengthOfString: Get string length
    Result: 13
    Duration: 2.5ms
    Executed at: 2024-01-15T10:30:01Z

🔄 Demonstrating database/sql compatibility...
  • Using standard database/sql interface:
    Compatibility: Working
  • Mixed usage (extended + standard):
    Extended: [Extended client works]
    Standard: Standard interface works

🎉 Extended client demonstration completed!
```

## 🆚 Comparison: Old vs New

### Old Approach (database/sql only)
```go
// System command - feels forced
rows, err := db.Query("COMMAND:df -h")
// Manual parsing of results...

// Function call - complex JSON handling
funcReq := FunctionRequest{...}
jsonData, _ := json.Marshal(funcReq)
rows, err := db.Query("FUNCTION:" + string(jsonData))
// Manual JSON parsing...
```

### New Approach (Extended Client)
```go
// System command - clean and typed
result, err := bc.ExecCommand("df -h")
fmt.Printf("Exit code: %d\n", result.ExitCode)
fmt.Printf("Output: %v\n", result.Stdout)

// Function call - simple and type-safe
result, err := bc.ExecFunction("myFunction", 
    client.StringParam("hello"),
    client.IntParam(42),
)
fmt.Printf("Result: %v\n", result.Result)
```

## 🎯 Benefits

1. **Cleaner Interface**: Separate methods for different operation types
2. **Type Safety**: Structured results with proper Go types
3. **Better Error Handling**: Comprehensive error information
4. **Metadata**: Execution time, timestamps, and performance data
5. **Backward Compatibility**: Full access to underlying `sql.DB`
6. **Easier Testing**: Mockable interface for unit tests
7. **Better Documentation**: Clear method signatures and purpose

## 🔧 Advanced Usage

### Custom Parameter Types
```go
// Create custom parameter
param := client.FunctionParam{
    Type:  "custom",
    Value: customStruct{Field: "value"},
}

result, err := bc.ExecFunction("customFunction", param)
```

### Transaction Support
```go
// Start transaction
tx, err := bc.Begin()
if err != nil {
    log.Fatal(err)
}

// Use transaction...
err = tx.Commit()
```

### Prepared Statements
```go
// Prepare statement
stmt, err := bc.Prepare("SELECT * FROM users WHERE id = ?")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

// Execute prepared statement
rows, err := stmt.Query(123)
```

This extended client provides a much cleaner and more intuitive interface while maintaining full compatibility with existing code that uses the standard `database/sql` interface.