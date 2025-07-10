package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lordbasex/burrowctl/client"
)

func main() {
	// Create DSN for connection
	dsn := "deviceID=extended-client-demo&amqp_uri=amqp://guest:guest@localhost:5672/&timeout=30s&debug=true"

	// Create extended burrow client
	bc, err := client.NewBurrowClient(dsn)
	if err != nil {
		log.Fatal("Failed to create burrow client:", err)
	}
	defer bc.Close()

	// Test connection
	fmt.Println("🔌 Testing connection...")
	if err := bc.Ping(); err != nil {
		log.Fatal("Connection test failed:", err)
	}
	fmt.Println("✅ Connection successful!")

	// Demonstrate SQL queries using the extended client
	fmt.Println("\n📊 Executing SQL queries...")
	demonstrateSQL(bc)

	// Demonstrate system commands using the extended client
	fmt.Println("\n🖥️ Executing system commands...")
	demonstrateCommands(bc)

	// Demonstrate function calls using the extended client
	fmt.Println("\n⚙️ Executing custom functions...")
	demonstrateFunctions(bc)

	// Show compatibility with standard database/sql interface
	fmt.Println("\n🔄 Demonstrating database/sql compatibility...")
	demonstrateCompatibility(bc)

	fmt.Println("\n🎉 Extended client demonstration completed!")
}

// demonstrateSQL shows how to use the extended client for SQL operations
func demonstrateSQL(bc *client.BurrowClient) {
	// Simple query
	fmt.Println("  • Simple SELECT query:")
	rows, err := bc.Query("SELECT 'Hello' as message, 42 as number")
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var message string
		var number int
		if err := rows.Scan(&message, &number); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}
		fmt.Printf("    Result: %s, %d\n", message, number)
	}

	// Query with parameters
	fmt.Println("  • Parameterized query:")
	row := bc.QueryRow("SELECT ? as name, ? as age", "Alice", 30)
	var name string
	var age int
	if err := row.Scan(&name, &age); err != nil {
		log.Printf("QueryRow failed: %v", err)
		return
	}
	fmt.Printf("    Result: %s is %d years old\n", name, age)
}

// demonstrateCommands shows how to use the extended client for system commands
func demonstrateCommands(bc *client.BurrowClient) {
	commands := []string{
		"echo 'Hello from system command'",
		"date",
		"uname -a",
		"df -h",
		"ps aux | head -5",
	}

	for _, cmd := range commands {
		fmt.Printf("  • Executing: %s\n", cmd)
		result, err := bc.ExecCommand(cmd)
		if err != nil {
			log.Printf("Command failed: %v", err)
			continue
		}

		fmt.Printf("    Exit code: %d\n", result.ExitCode)
		fmt.Printf("    Executed at: %s\n", result.ExecutedAt.Format(time.RFC3339))
		
		if len(result.Stdout) > 0 {
			fmt.Printf("    Stdout:\n")
			for _, line := range result.Stdout {
				fmt.Printf("      %s\n", line)
			}
		}
		
		if len(result.Stderr) > 0 {
			fmt.Printf("    Stderr:\n")
			for _, line := range result.Stderr {
				fmt.Printf("      %s\n", line)
			}
		}
		fmt.Println()
	}
}

// demonstrateFunctions shows how to use the extended client for function calls
func demonstrateFunctions(bc *client.BurrowClient) {
	// Test various built-in functions
	functions := []struct {
		name   string
		params []client.FunctionParam
		desc   string
	}{
		{
			name:   "lengthOfString",
			params: []client.FunctionParam{client.StringParam("Hello, World!")},
			desc:   "Get string length",
		},
		{
			name:   "addIntegers",
			params: []client.FunctionParam{client.IntParam(15), client.IntParam(27)},
			desc:   "Add two integers",
		},
		{
			name:   "getCurrentTimestamp",
			params: []client.FunctionParam{},
			desc:   "Get current timestamp",
		},
		{
			name:   "generateUUID",
			params: []client.FunctionParam{},
			desc:   "Generate UUID",
		},
		{
			name:   "encodeBase64",
			params: []client.FunctionParam{client.StringParam("Hello, Base64!")},
			desc:   "Encode string to Base64",
		},
		{
			name:   "calculateHash",
			params: []client.FunctionParam{client.StringParam("test data")},
			desc:   "Calculate SHA256 hash",
		},
		{
			name:   "validateEmail",
			params: []client.FunctionParam{client.StringParam("user@example.com")},
			desc:   "Validate email address",
		},
		{
			name:   "generateRandomString",
			params: []client.FunctionParam{client.IntParam(10)},
			desc:   "Generate random string",
		},
	}

	for _, fn := range functions {
		fmt.Printf("  • %s: %s\n", fn.name, fn.desc)
		result, err := bc.ExecFunction(fn.name, fn.params...)
		if err != nil {
			log.Printf("Function failed: %v", err)
			continue
		}

		fmt.Printf("    Result: %v\n", result.Result)
		fmt.Printf("    Duration: %s\n", result.Duration)
		fmt.Printf("    Executed at: %s\n", result.ExecutedAt.Format(time.RFC3339))
		
		if result.Error != "" {
			fmt.Printf("    Error: %s\n", result.Error)
		}
		fmt.Println()
	}

	// Test complex function with JSON parameter
	fmt.Println("  • parseJSON: Parse JSON string")
	jsonStr := `{"name": "John", "age": 30, "city": "New York"}`
	result, err := bc.ExecFunction("parseJSON", client.StringParam(jsonStr))
	if err != nil {
		log.Printf("parseJSON failed: %v", err)
		return
	}
	fmt.Printf("    Result: %v\n", result.Result)
	fmt.Printf("    Duration: %s\n", result.Duration)
}

// demonstrateCompatibility shows that the extended client maintains compatibility with database/sql
func demonstrateCompatibility(bc *client.BurrowClient) {
	// Get the underlying sql.DB instance
	db := bc.DB()

	// Use standard database/sql methods
	fmt.Println("  • Using standard database/sql interface:")
	rows, err := db.Query("SELECT 'Compatibility' as feature, 'Working' as status")
	if err != nil {
		log.Printf("Standard query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var feature, status string
		if err := rows.Scan(&feature, &status); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}
		fmt.Printf("    %s: %s\n", feature, status)
	}

	// Show that both interfaces work on the same connection
	fmt.Println("  • Mixed usage (extended + standard):")
	
	// Use extended client
	cmdResult, err := bc.ExecCommand("echo 'Extended client works'")
	if err != nil {
		log.Printf("Extended command failed: %v", err)
		return
	}
	fmt.Printf("    Extended: %v\n", cmdResult.Stdout)

	// Use standard interface
	row := db.QueryRow("SELECT 'Standard interface works' as message")
	var message string
	if err := row.Scan(&message); err != nil {
		log.Printf("Standard query failed: %v", err)
		return
	}
	fmt.Printf("    Standard: %s\n", message)
}