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

	// Test basic transaction functionality
	fmt.Println("🔄 Testing Basic Transaction Functionality")
	testBasicTransaction(db)

	// Test transaction rollback
	fmt.Println("\n🔄 Testing Transaction Rollback")
	testTransactionRollback(db)

	// Test transaction timeout
	fmt.Println("\n🔄 Testing Transaction Timeout")
	testTransactionTimeout(db)

	// Test multiple operations in transaction
	fmt.Println("\n🔄 Testing Multiple Operations in Transaction")
	testMultipleOperations(db)

	fmt.Println("\n✅ All transaction tests completed!")
}

// testBasicTransaction demonstrates basic transaction commit functionality
func testBasicTransaction(db *sql.DB) {
	fmt.Println("→ Starting basic transaction test...")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return
	}

	// Execute some operations within transaction
	_, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
	if err != nil {
		log.Printf("❌ Failed to insert user: %v", err)
		tx.Rollback()
		return
	}

	// Query within transaction
	rows, err := tx.Query("SELECT name, email FROM users WHERE email = ?", "john@example.com")
	if err != nil {
		log.Printf("❌ Failed to query user: %v", err)
		tx.Rollback()
		return
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var name, email string
		if err := rows.Scan(&name, &email); err != nil {
			log.Printf("❌ Failed to scan row: %v", err)
			continue
		}
		fmt.Printf("  📄 Found user: %s (%s)\n", name, email)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("❌ Failed to commit transaction: %v", err)
		return
	}

	fmt.Println("✅ Basic transaction test completed successfully!")
}

// testTransactionRollback demonstrates transaction rollback functionality
func testTransactionRollback(db *sql.DB) {
	fmt.Println("→ Starting transaction rollback test...")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return
	}

	// Execute some operations
	_, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Jane Doe", "jane@example.com")
	if err != nil {
		log.Printf("❌ Failed to insert user: %v", err)
		tx.Rollback()
		return
	}

	fmt.Println("  📄 User inserted within transaction")

	// Intentionally rollback
	if err := tx.Rollback(); err != nil {
		log.Printf("❌ Failed to rollback transaction: %v", err)
		return
	}

	fmt.Println("  🔄 Transaction rolled back successfully")

	// Verify the user was not actually inserted
	rows, err := db.Query("SELECT name FROM users WHERE email = ?", "jane@example.com")
	if err != nil {
		log.Printf("❌ Failed to query user: %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		fmt.Println("❌ User found after rollback (this shouldn't happen)")
	} else {
		fmt.Println("✅ User correctly not found after rollback")
	}

	fmt.Println("✅ Transaction rollback test completed successfully!")
}

// testTransactionTimeout demonstrates transaction timeout handling
func testTransactionTimeout(db *sql.DB) {
	fmt.Println("→ Starting transaction timeout test...")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return
	}

	// Execute quick operation
	_, err = tx.Exec("SELECT 1")
	if err != nil {
		log.Printf("❌ Failed to execute query: %v", err)
		tx.Rollback()
		return
	}

	fmt.Println("  📄 Quick operation completed")

	// Wait a bit to test timeout handling
	time.Sleep(2 * time.Second)

	// Try to commit
	if err := tx.Commit(); err != nil {
		log.Printf("⚠️  Transaction commit failed (expected for timeout test): %v", err)
	} else {
		fmt.Println("✅ Transaction committed successfully")
	}

	fmt.Println("✅ Transaction timeout test completed!")
}

// testMultipleOperations demonstrates multiple operations in a single transaction
func testMultipleOperations(db *sql.DB) {
	fmt.Println("→ Starting multiple operations test...")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return
	}

	// Operation 1: Insert user
	_, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Alice Smith", "alice@example.com")
	if err != nil {
		log.Printf("❌ Failed to insert user: %v", err)
		tx.Rollback()
		return
	}
	fmt.Println("  📄 User inserted")

	// Operation 2: Update user
	_, err = tx.Exec("UPDATE users SET name = ? WHERE email = ?", "Alice Johnson", "alice@example.com")
	if err != nil {
		log.Printf("❌ Failed to update user: %v", err)
		tx.Rollback()
		return
	}
	fmt.Println("  📄 User updated")

	// Operation 3: Query user
	rows, err := tx.Query("SELECT name, email FROM users WHERE email = ?", "alice@example.com")
	if err != nil {
		log.Printf("❌ Failed to query user: %v", err)
		tx.Rollback()
		return
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var name, email string
		if err := rows.Scan(&name, &email); err != nil {
			log.Printf("❌ Failed to scan row: %v", err)
			continue
		}
		fmt.Printf("  📄 Final user state: %s (%s)\n", name, email)
	}

	// Commit all operations
	if err := tx.Commit(); err != nil {
		log.Printf("❌ Failed to commit transaction: %v", err)
		return
	}

	fmt.Println("✅ Multiple operations test completed successfully!")
}