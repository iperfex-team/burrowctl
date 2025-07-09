package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

	fmt.Println("🛡️ SQL Validation Security Test")
	fmt.Println("===============================")

	// Test 1: Valid queries (should pass)
	fmt.Println("\n→ Test 1: Valid Queries (should pass)")
	testValidQueries(db)

	// Test 2: SQL Injection attempts (should be blocked)
	fmt.Println("\n→ Test 2: SQL Injection Attempts (should be blocked)")
	testSQLInjectionAttempts(db)

	// Test 3: Command policy violations (should be blocked)
	fmt.Println("\n→ Test 3: Command Policy Violations (should be blocked)")
	testCommandViolations(db)

	// Test 4: Structure violations (should be detected)
	fmt.Println("\n→ Test 4: Structure Violations (should be detected)")
	testStructureViolations(db)

	// Test 5: Parameter validation (should show warnings)
	fmt.Println("\n→ Test 5: Parameter Validation (should show warnings)")
	testParameterValidation(db)

	fmt.Println("\n✅ SQL validation test completed!")
	fmt.Println("📊 Check server logs for detailed validation results and statistics.")
}

// testValidQueries tests legitimate queries that should pass validation
func testValidQueries(db *sql.DB) {
	validQueries := []struct {
		query  string
		params []interface{}
		desc   string
	}{
		{
			query:  "SELECT id, name FROM users WHERE active = ?",
			params: []interface{}{true},
			desc:   "Basic SELECT with parameter",
		},
		{
			query:  "SELECT COUNT(*) FROM users",
			params: []interface{}{},
			desc:   "COUNT query",
		},
		{
			query:  "SELECT * FROM products WHERE price BETWEEN ? AND ?",
			params: []interface{}{10.0, 100.0},
			desc:   "BETWEEN clause with parameters",
		},
		{
			query:  "INSERT INTO logs (message, timestamp) VALUES (?, NOW())",
			params: []interface{}{"Test message"},
			desc:   "INSERT statement",
		},
	}

	for i, test := range validQueries {
		fmt.Printf("  %d. %s\n", i+1, test.desc)
		fmt.Printf("     Query: %s\n", test.query)

		_, err := db.Query(test.query, test.params...)
		if err != nil {
			fmt.Printf("     ❌ Failed: %v\n", err)
		} else {
			fmt.Printf("     ✅ Passed validation\n")
		}
	}
}

// testSQLInjectionAttempts tests various SQL injection techniques
func testSQLInjectionAttempts(db *sql.DB) {
	injectionAttempts := []struct {
		query string
		desc  string
	}{
		{
			query: "SELECT * FROM users WHERE id = 1 OR 1=1",
			desc:  "Boolean-based injection (OR 1=1)",
		},
		{
			query: "SELECT * FROM users WHERE name = 'admin' UNION SELECT password FROM admin_users",
			desc:  "UNION-based injection",
		},
		{
			query: "SELECT * FROM users WHERE id = 1; DROP TABLE users;",
			desc:  "Stacked queries (DROP TABLE)",
		},
		{
			query: "SELECT * FROM users WHERE name = 'test' AND SLEEP(5)",
			desc:  "Time-based injection (SLEEP)",
		},
		{
			query: "SELECT * FROM users WHERE id = 1 AND (SELECT COUNT(*) FROM information_schema.tables) > 0",
			desc:  "Information schema attack",
		},
		{
			query: "SELECT * FROM users WHERE name = CHAR(97,100,109,105,110)",
			desc:  "CHAR encoding injection",
		},
		{
			query: "SELECT * FROM users WHERE id = 1 /* comment */ OR /* comment */ 1=1",
			desc:  "Comment-based injection",
		},
	}

	for i, test := range injectionAttempts {
		fmt.Printf("  %d. %s\n", i+1, test.desc)
		fmt.Printf("     Query: %s\n", truncateString(test.query, 80))

		_, err := db.Query(test.query)
		if err != nil {
			if containsValidationError(err.Error()) {
				fmt.Printf("     ✅ Blocked by validation: %s\n", extractValidationError(err.Error()))
			} else {
				fmt.Printf("     ❌ Unexpected error: %v\n", err)
			}
		} else {
			fmt.Printf("     ⚠️  WARNING: Injection attempt was not blocked!\n")
		}
	}
}

// testCommandViolations tests queries that violate command policies
func testCommandViolations(db *sql.DB) {
	violationQueries := []struct {
		query string
		desc  string
	}{
		{
			query: "DROP TABLE users",
			desc:  "DROP TABLE command",
		},
		{
			query: "CREATE TABLE test_table (id INT)",
			desc:  "CREATE TABLE command",
		},
		{
			query: "ALTER TABLE users ADD COLUMN test VARCHAR(255)",
			desc:  "ALTER TABLE command",
		},
		{
			query: "TRUNCATE TABLE logs",
			desc:  "TRUNCATE command",
		},
		{
			query: "GRANT ALL PRIVILEGES ON *.* TO 'user'@'%'",
			desc:  "GRANT privileges command",
		},
	}

	for i, test := range violationQueries {
		fmt.Printf("  %d. %s\n", i+1, test.desc)
		fmt.Printf("     Query: %s\n", test.query)

		_, err := db.Query(test.query)
		if err != nil {
			if containsValidationError(err.Error()) {
				fmt.Printf("     ✅ Blocked by policy: %s\n", extractValidationError(err.Error()))
			} else {
				fmt.Printf("     ❌ Unexpected error: %v\n", err)
			}
		} else {
			fmt.Printf("     ⚠️  WARNING: Command violation was not blocked!\n")
		}
	}
}

// testStructureViolations tests queries with structural problems
func testStructureViolations(db *sql.DB) {
	structureQueries := []struct {
		query string
		desc  string
	}{
		{
			query: "SELECT * FROM users WHERE (name = 'test'",
			desc:  "Unbalanced parentheses",
		},
		{
			query: "SELECT * FROM users WHERE name = 'test",
			desc:  "Unbalanced quotes",
		},
		{
			query: "SELECT * FROM users WHERE name = 'test'; SELECT * FROM admin",
			desc:  "Multiple statements",
		},
		{
			query: "SELECT * FROM users /* unclosed comment",
			desc:  "Unclosed comment block",
		},
	}

	for i, test := range structureQueries {
		fmt.Printf("  %d. %s\n", i+1, test.desc)
		fmt.Printf("     Query: %s\n", test.query)

		_, err := db.Query(test.query)
		if err != nil {
			if containsValidationError(err.Error()) {
				fmt.Printf("     ✅ Detected structural issue: %s\n", extractValidationError(err.Error()))
			} else {
				fmt.Printf("     ❌ Unexpected error: %v\n", err)
			}
		} else {
			fmt.Printf("     ⚠️  Structural issue not detected\n")
		}
	}
}

// testParameterValidation tests parameter validation
func testParameterValidation(db *sql.DB) {
	paramTests := []struct {
		query  string
		params []interface{}
		desc   string
	}{
		{
			query:  "SELECT * FROM users WHERE name = ?",
			params: []interface{}{"admin' OR '1'='1"},
			desc:   "Parameter with SQL injection attempt",
		},
		{
			query:  "SELECT * FROM users WHERE description = ?",
			params: []interface{}{"Normal description"},
			desc:   "Normal parameter",
		},
		{
			query:  "SELECT * FROM users WHERE query_text = ?",
			params: []interface{}{"SELECT * FROM secret_table"},
			desc:   "Parameter containing SQL keywords",
		},
	}

	for i, test := range paramTests {
		fmt.Printf("  %d. %s\n", i+1, test.desc)
		fmt.Printf("     Query: %s\n", test.query)
		fmt.Printf("     Params: %v\n", test.params)

		_, err := db.Query(test.query, test.params...)
		if err != nil {
			fmt.Printf("     Result: %v\n", err)
		} else {
			fmt.Printf("     ✅ Query executed (check logs for warnings)\n")
		}
	}
}

// Helper functions
func containsValidationError(errMsg string) bool {
	return strings.Contains(errMsg, "SQL validation failed") || 
		   strings.Contains(errMsg, "validation")
}

func extractValidationError(errMsg string) string {
	if idx := strings.Index(errMsg, "SQL validation failed: "); idx != -1 {
		return errMsg[idx+len("SQL validation failed: "):]
	}
	return errMsg
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}