package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// Configuration
	dsn := "deviceID=my-device&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=30s&debug=true"

	// Connect to the database
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer db.Close()

	fmt.Println("🔄 Query Cache Performance Test")
	fmt.Println("================================")

	// Test 1: First execution (should be cache miss)
	fmt.Println("\n→ Test 1: First query execution (cache miss expected)")
	runQueryBenchmark(db, "SELECT * FROM users LIMIT 5", 1)

	// Test 2: Immediate re-execution (should be cache hit)
	fmt.Println("\n→ Test 2: Immediate re-execution (cache hit expected)")
	runQueryBenchmark(db, "SELECT * FROM users LIMIT 5", 1)

	// Test 3: Multiple executions of the same query
	fmt.Println("\n→ Test 3: Multiple executions of same query (should show cache hits)")
	runQueryBenchmark(db, "SELECT * FROM users LIMIT 5", 5)

	// Test 4: Different query (should be cache miss)
	fmt.Println("\n→ Test 4: Different query (cache miss expected)")
	runQueryBenchmark(db, "SELECT COUNT(*) FROM users", 1)

	// Test 5: Parameterized queries
	fmt.Println("\n→ Test 5: Parameterized queries")
	testParameterizedQueries(db)

	// Test 6: Write operations (should not be cached)
	fmt.Println("\n→ Test 6: Write operations (should not be cached)")
	testWriteOperations(db)

	// Test 7: Performance comparison
	fmt.Println("\n→ Test 7: Performance comparison (cached vs non-cached)")
	performanceComparison(db)

	fmt.Println("\n✅ Cache performance test completed!")
}

// runQueryBenchmark executes a query multiple times and measures performance
func runQueryBenchmark(db *sql.DB, query string, iterations int) {
	fmt.Printf("  Executing query %d time(s): %s\n", iterations, truncateString(query, 50))
	
	var totalDuration time.Duration
	var results [][]interface{}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("❌ Query failed: %v", err)
			continue
		}

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			log.Printf("❌ Failed to get columns: %v", err)
			rows.Close()
			continue
		}

		// Process results
		var rowData [][]interface{}
		for rows.Next() {
			// Create destinations for scanning
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}

			// Scan row
			if err := rows.Scan(pointers...); err != nil {
				log.Printf("❌ Failed to scan row: %v", err)
				continue
			}

			rowData = append(rowData, values)
		}
		rows.Close()

		duration := time.Since(start)
		totalDuration += duration

		// Store results for first iteration
		if i == 0 {
			results = rowData
		}

		fmt.Printf("    Iteration %d: %v (%d rows)\n", i+1, duration, len(rowData))
	}

	avgDuration := totalDuration / time.Duration(iterations)
	fmt.Printf("  📊 Average execution time: %v\n", avgDuration)
	fmt.Printf("  📄 Total rows returned: %d\n", len(results))
}

// testParameterizedQueries tests caching with different parameters
func testParameterizedQueries(db *sql.DB) {
	queries := []struct {
		query  string
		params []interface{}
		name   string
	}{
		{"SELECT * FROM users WHERE id = ?", []interface{}{1}, "User ID 1"},
		{"SELECT * FROM users WHERE id = ?", []interface{}{1}, "User ID 1 (repeat)"},
		{"SELECT * FROM users WHERE id = ?", []interface{}{2}, "User ID 2"},
		{"SELECT * FROM users WHERE id = ?", []interface{}{1}, "User ID 1 (again)"},
	}

	for _, test := range queries {
		fmt.Printf("  Testing: %s\n", test.name)
		start := time.Now()

		rows, err := db.Query(test.query, test.params...)
		if err != nil {
			log.Printf("❌ Query failed: %v", err)
			continue
		}

		// Count rows
		rowCount := 0
		for rows.Next() {
			rowCount++
		}
		rows.Close()

		duration := time.Since(start)
		fmt.Printf("    Result: %v (%d rows)\n", duration, rowCount)
	}
}

// testWriteOperations tests that write operations are not cached
func testWriteOperations(db *sql.DB) {
	writeQueries := []string{
		"INSERT INTO users (name, email) VALUES ('Cache Test', 'cache@test.com')",
		"UPDATE users SET name = 'Updated Cache Test' WHERE email = 'cache@test.com'",
		"DELETE FROM users WHERE email = 'cache@test.com'",
	}

	for _, query := range writeQueries {
		fmt.Printf("  Executing: %s\n", truncateString(query, 60))
		start := time.Now()

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("❌ Write operation failed: %v", err)
			continue
		}

		duration := time.Since(start)
		fmt.Printf("    Completed: %v (should not be cached)\n", duration)
	}
}

// performanceComparison compares performance of repeated queries
func performanceComparison(db *sql.DB) {
	query := "SELECT * FROM users ORDER BY id LIMIT 10"
	iterations := 10

	fmt.Printf("  Running %d iterations of: %s\n", iterations, truncateString(query, 50))

	times := make([]time.Duration, iterations)
	for i := 0; i < iterations; i++ {
		start := time.Now()

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("❌ Query failed: %v", err)
			continue
		}

		// Process results
		rowCount := 0
		for rows.Next() {
			rowCount++
		}
		rows.Close()

		times[i] = time.Since(start)
		fmt.Printf("    Iteration %d: %v\n", i+1, times[i])
	}

	// Calculate statistics
	if len(times) > 0 {
		var total time.Duration
		min := times[0]
		max := times[0]

		for _, t := range times {
			total += t
			if t < min {
				min = t
			}
			if t > max {
				max = t
			}
		}

		avg := total / time.Duration(len(times))
		fmt.Printf("\n  📊 Performance Summary:\n")
		fmt.Printf("    Average: %v\n", avg)
		fmt.Printf("    Minimum: %v (likely cache hit)\n", min)
		fmt.Printf("    Maximum: %v (likely cache miss)\n", max)
		fmt.Printf("    Speedup: %.2fx (max/min)\n", float64(max)/float64(min))
	}
}

// truncateString truncates a string for display purposes
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}