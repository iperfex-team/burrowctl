package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// Configuración de flags para demostrar las nuevas características
	var (
		deviceID        = flag.String("device", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb", "Device ID")
		amqpURL         = flag.String("amqp", "amqp://burrowuser:burrowpass123@localhost:5672/", "RabbitMQ URL")
		timeout         = flag.Duration("timeout", 10*time.Second, "Query timeout (e.g., 5s, 30s, 1m)")
		debug           = flag.Bool("debug", true, "Enable debug logging")
		usePrepared     = flag.Bool("prepared", false, "Use prepared statements")
		stressTest      = flag.Bool("stress", false, "Run stress test to demonstrate rate limiting")
		concurrent      = flag.Int("concurrent", 5, "Number of concurrent connections for stress test")
		requests        = flag.Int("requests", 50, "Number of requests per connection for stress test")
		reconnectDemo   = flag.Bool("reconnect-demo", false, "Demonstrate automatic reconnection")
		showHelp        = flag.Bool("help", false, "Show this help message")
	)
	flag.Parse()

	if *showHelp {
		showAdvancedHelp()
		return
	}

	// Construir DSN con configuración personalizada
	dsn := fmt.Sprintf("deviceID=%s&amqp_uri=%s&timeout=%s&debug=%t",
		*deviceID, *amqpURL, timeout.String(), *debug)

	fmt.Printf("🗃️  Advanced burrowctl SQL Example\n")
	fmt.Printf("================================================\n")
	fmt.Printf("📱 Device ID: %s\n", *deviceID)
	fmt.Printf("🐰 RabbitMQ: %s\n", *amqpURL)
	fmt.Printf("⏱️  Timeout: %s\n", timeout.String())
	fmt.Printf("🐛 Debug: %t\n", *debug)
	fmt.Printf("📝 Prepared Statements: %t\n", *usePrepared)
	fmt.Printf("📡 DSN: %s\n", dsn)
	fmt.Println()

	// Diferentes modos de demostración
	switch {
	case *stressTest:
		runStressTest(dsn, *concurrent, *requests)
	case *reconnectDemo:
		runReconnectDemo(dsn)
	case *usePrepared:
		runPreparedStatementsDemo(dsn)
	default:
		runBasicDemo(dsn)
	}
}

func showAdvancedHelp() {
	fmt.Println("🚀 Advanced burrowctl SQL Example")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("This example demonstrates the new enterprise features:")
	fmt.Println("• 🔄 Automatic Reconnection - Handles connection failures gracefully")
	fmt.Println("• ⚡ Prepared Statements - Better performance and security")
	fmt.Println("• 🏗️  Worker Pool - Concurrent processing on server")
	fmt.Println("• 🛡️  Rate Limiting - Protection against abuse")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run advanced-main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -device string     Device ID (default: SHA256 hash)")
	fmt.Println("  -amqp string       RabbitMQ URL (default: localhost)")
	fmt.Println("  -timeout duration  Query timeout (default: 10s)")
	fmt.Println("  -debug            Enable debug logging (default: true)")
	fmt.Println("  -prepared         Use prepared statements demo")
	fmt.Println("  -stress           Run stress test (rate limiting demo)")
	fmt.Println("  -concurrent int   Concurrent connections for stress test (default: 5)")
	fmt.Println("  -requests int     Requests per connection (default: 50)")
	fmt.Println("  -reconnect-demo   Demonstrate automatic reconnection")
	fmt.Println("  -help             Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Basic usage with custom timeout")
	fmt.Println("  go run advanced-main.go -timeout=30s")
	fmt.Println()
	fmt.Println("  # Prepared statements demo")
	fmt.Println("  go run advanced-main.go -prepared")
	fmt.Println()
	fmt.Println("  # Stress test to trigger rate limiting")
	fmt.Println("  go run advanced-main.go -stress -concurrent=10 -requests=100")
	fmt.Println()
	fmt.Println("  # Reconnection demo")
	fmt.Println("  go run advanced-main.go -reconnect-demo")
	fmt.Println()
	fmt.Println("  # Custom configuration")
	fmt.Println("  go run advanced-main.go -device=mydevice -amqp=amqp://user:pass@host:5672/ -timeout=1m")
}

func runBasicDemo(dsn string) {
	fmt.Println("🎯 Running Basic Demo")
	fmt.Println("---------------------")

	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("❌ Error connecting:", err)
	}
	defer db.Close()

	// Test query from command line or use default
	query := "SELECT 'Hello' as greeting, 'World' as target, NOW() as timestamp"
	if len(flag.Args()) > 0 {
		query = flag.Args()[0]
	}

	fmt.Printf("📊 Executing: %s\n", query)

	start := time.Now()
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("❌ Error executing query:", err)
	}
	defer rows.Close()

	printResults(rows)
	fmt.Printf("⏱️  Query completed in: %v\n", time.Since(start))
	fmt.Println()
	fmt.Println("✅ Basic demo completed!")
	fmt.Println("💡 Try: go run advanced-main.go -prepared")
}

func runPreparedStatementsDemo(dsn string) {
	fmt.Println("🎯 Running Prepared Statements Demo")
	fmt.Println("-----------------------------------")

	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("❌ Error connecting:", err)
	}
	defer db.Close()

	// Prepare a statement with parameters
	query := "SELECT ? as message, ? as number, ? as flag"
	fmt.Printf("📝 Preparing statement: %s\n", query)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("❌ Error preparing statement:", err)
	}
	defer stmt.Close()

	// Execute with different parameters multiple times
	testData := [][]interface{}{
		{"Hello World", 42, true},
		{"Prepared Statement", 100, false},
		{"Performance Test", 999, true},
	}

	for i, params := range testData {
		fmt.Printf("\n📋 Execution %d with params: %v\n", i+1, params)
		
		start := time.Now()
		rows, err := stmt.Query(params...)
		if err != nil {
			log.Printf("❌ Error executing prepared statement: %v", err)
			continue
		}

		printResults(rows)
		rows.Close()
		
		fmt.Printf("⏱️  Execution %d completed in: %v\n", i+1, time.Since(start))
	}

	fmt.Println()
	fmt.Println("✅ Prepared statements demo completed!")
	fmt.Println("💡 Benefits: Better performance, SQL injection protection")
	fmt.Println("💡 Try: go run advanced-main.go -stress")
}

func runStressTest(dsn string, concurrent int, requests int) {
	fmt.Println("🎯 Running Stress Test (Rate Limiting Demo)")
	fmt.Println("-------------------------------------------")
	fmt.Printf("🏗️  Concurrent connections: %d\n", concurrent)
	fmt.Printf("📊 Requests per connection: %d\n", requests)
	fmt.Printf("📈 Total requests: %d\n", concurrent*requests)
	fmt.Println()

	resultChan := make(chan TestResult, concurrent)

	// Start concurrent workers
	for i := 0; i < concurrent; i++ {
		go stressWorker(dsn, i+1, requests, resultChan)
	}

	// Collect results
	var totalRequests, totalErrors, rateLimited int
	var totalDuration time.Duration

	fmt.Println("📊 Live Results:")
	fmt.Println("Worker | Success | Errors | Rate Limited | Avg Time")
	fmt.Println("-------|---------|--------|--------------|----------")

	for i := 0; i < concurrent; i++ {
		result := <-resultChan
		totalRequests += result.TotalRequests
		totalErrors += result.Errors
		rateLimited += result.RateLimited
		totalDuration += result.TotalDuration

		avgTime := time.Duration(0)
		if result.TotalRequests > 0 {
			avgTime = result.TotalDuration / time.Duration(result.TotalRequests)
		}

		fmt.Printf("   %-3d | %-7d | %-6d | %-12d | %8v\n",
			result.WorkerID, 
			result.TotalRequests-result.Errors, 
			result.Errors,
			result.RateLimited,
			avgTime)
	}

	fmt.Println()
	fmt.Println("📈 Final Statistics:")
	fmt.Printf("✅ Total Successful: %d\n", totalRequests-totalErrors)
	fmt.Printf("❌ Total Errors: %d\n", totalErrors)
	fmt.Printf("🛡️  Rate Limited: %d\n", rateLimited)
	fmt.Printf("⏱️  Average Response Time: %v\n", totalDuration/time.Duration(totalRequests))
	
	if rateLimited > 0 {
		fmt.Println()
		fmt.Println("🎯 Rate Limiting Demonstration Successful!")
		fmt.Printf("   Server protected against %d excessive requests\n", rateLimited)
		fmt.Println("   💡 This shows the rate limiter is working correctly")
	}

	fmt.Println()
	fmt.Println("✅ Stress test completed!")
	fmt.Println("💡 Try: go run advanced-main.go -reconnect-demo")
}

type TestResult struct {
	WorkerID      int
	TotalRequests int
	Errors        int
	RateLimited   int
	TotalDuration time.Duration
}

func stressWorker(dsn string, workerID, requests int, resultChan chan<- TestResult) {
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Printf("Worker %d: Failed to connect: %v", workerID, err)
		resultChan <- TestResult{WorkerID: workerID, Errors: requests}
		return
	}
	defer db.Close()

	var errors, rateLimited int
	var totalDuration time.Duration

	for i := 0; i < requests; i++ {
		start := time.Now()
		
		// Update query with current request number
		currentQuery := fmt.Sprintf("SELECT 'Worker %d' as worker, %d as request_num", workerID, i+1)
		
		rows, err := db.Query(currentQuery)
		duration := time.Since(start)
		totalDuration += duration

		if err != nil {
			if fmt.Sprintf("%v", err) == "server error: Rate limit exceeded. Please slow down your requests." {
				rateLimited++
			}
			errors++
		} else {
			rows.Close()
		}

		// Small delay to simulate realistic usage
		time.Sleep(50 * time.Millisecond)
	}

	resultChan <- TestResult{
		WorkerID:      workerID,
		TotalRequests: requests,
		Errors:        errors,
		RateLimited:   rateLimited,
		TotalDuration: totalDuration,
	}
}

func runReconnectDemo(dsn string) {
	fmt.Println("🎯 Running Automatic Reconnection Demo")
	fmt.Println("--------------------------------------")

	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("❌ Error connecting:", err)
	}
	defer db.Close()

	fmt.Println("🔄 This demo shows automatic reconnection handling")
	fmt.Println("   The client will automatically reconnect if connection is lost")
	fmt.Println()

	// Run queries with intervals to demonstrate persistent connection
	for i := 1; i <= 10; i++ {
		fmt.Printf("📊 Query %d/10: Testing connection...\n", i)
		
		start := time.Now()
		testQuery := fmt.Sprintf("SELECT 'Reconnection Test' as test, %d as iteration, NOW() as timestamp", i)
		
		rows, err := db.Query(testQuery)
		if err != nil {
			fmt.Printf("❌ Query %d failed: %v\n", i, err)
			fmt.Println("🔄 Client will attempt automatic reconnection...")
		} else {
			printResults(rows)
			rows.Close()
			fmt.Printf("✅ Query %d completed in: %v\n", i, time.Since(start))
		}

		if i == 5 {
			fmt.Println()
			fmt.Println("💡 Simulate connection loss now (disconnect RabbitMQ) to see auto-reconnection")
			fmt.Println("   The client will automatically handle reconnection with exponential backoff")
			fmt.Println("   Press Ctrl+C to stop the demo")
			fmt.Println()
		}

		// Wait between queries
		time.Sleep(3 * time.Second)
	}

	fmt.Println()
	fmt.Println("✅ Reconnection demo completed!")
	fmt.Println("💡 Features demonstrated:")
	fmt.Println("   • Automatic reconnection with exponential backoff")
	fmt.Println("   • Connection health monitoring")
	fmt.Println("   • Transparent error handling")
}

func printResults(rows *sql.Rows) {
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	// Print headers
	for i, col := range columns {
		if i > 0 {
			fmt.Printf(" | ")
		}
		fmt.Printf("%-15s", col)
	}
	fmt.Println()

	// Print separator
	for i := range columns {
		if i > 0 {
			fmt.Printf("-+-")
		}
		fmt.Printf("%-15s", "---------------")
	}
	fmt.Println()

	// Print data
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		for i, val := range values {
			if i > 0 {
				fmt.Printf(" | ")
			}
			if val == nil {
				fmt.Printf("%-15s", "<NULL>")
			} else {
				// Convert byte arrays to strings for better display
				if b, ok := val.([]byte); ok {
					fmt.Printf("%-15s", string(b))
				} else {
					fmt.Printf("%-15v", val)
				}
			}
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating results: %v", err)
	}
}