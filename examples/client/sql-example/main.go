package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// DSN con credenciales hardcodeadas para RabbitMQ
	dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true"

	// Abrir conexión usando el driver rabbitsql
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer db.Close()

	// Permitir query personalizada desde argumentos
	query := "SELECT id, name FROM users"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Printf("🗃️ burrowctl SQL Example\n")
	fmt.Printf("📡 Device ID: %s\n", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	fmt.Printf("🚀 Executing query: %s\n", query)

	log.Printf("Executing query: %s", query)

	// Ejecutar query SQL directamente (sin prefijo)
	// El cliente detectará automáticamente que es SQL y enviará type="sql"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	fmt.Println("\n📊 --- SQL Results ---")

	// Obtener nombres de columnas
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting columns:", err)
	}

	// Imprimir encabezados
	for i, col := range columns {
		if i > 0 {
			fmt.Printf(" | ")
		}
		fmt.Printf("%-15s", col)
	}
	fmt.Println()

	// Imprimir línea separadora
	for i := range columns {
		if i > 0 {
			fmt.Printf("-+-")
		}
		fmt.Printf("%-15s", "---------------")
	}
	fmt.Println()

	// Procesar resultados dinámicamente
	for rows.Next() {
		// Crear slice para escanear valores
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.Fatal("Error scanning result:", err)
		}

		// Imprimir valores
		for i, val := range values {
			if i > 0 {
				fmt.Printf(" | ")
			}
			if val == nil {
				fmt.Printf("%-15s", "<NULL>")
			} else {
				fmt.Printf("%-15v", val)
			}
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating results:", err)
	}

	fmt.Println("\n✅ SQL query executed successfully!")
	fmt.Println("\n💡 Usage tips:")
	fmt.Println("   - Default query: SELECT id, name FROM users")
	fmt.Println("   - Custom query: go run main.go \"SELECT * FROM products\"")
	fmt.Println("   - The client processes SQL queries as type='sql' by default")
	fmt.Println("   - Use quotes for complex queries with spaces")
}
