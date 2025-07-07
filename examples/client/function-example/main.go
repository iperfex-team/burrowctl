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

	// Verificar argumentos de línea de comandos
	functionName := "processData"
	if len(os.Args) > 1 {
		functionName = os.Args[1]
	}

	log.Printf("Executing remote function: %s", functionName)
	fmt.Printf("📡 Connecting to device: %s\n", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")

	// Ejecutar función remota usando el prefijo "FUNCTION:"
	// El cliente detectará automáticamente que es una función y enviará type="function"
	rows, err := db.Query("FUNCTION:" + functionName)
	if err != nil {
		log.Fatal("Error executing function:", err)
	}
	defer rows.Close()

	fmt.Println("\n🔧 --- Function Results ---")
	fmt.Printf("%-40s\n", "Message")
	fmt.Println("--------------------------------------------")

	// Procesar resultados
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			log.Fatal("Error scanning result:", err)
		}
		fmt.Printf("%-40s\n", message)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating results:", err)
	}

	fmt.Println("\n✅ Function executed successfully")
	fmt.Println("\n💡 Tips:")
	fmt.Println("   - You can specify a function name as argument:")
	fmt.Println("     go run main.go myCustomFunction")
	fmt.Println("   - The client automatically detects FUNCTION: prefix")
	fmt.Println("   - Server will process this as type='function'")
}
