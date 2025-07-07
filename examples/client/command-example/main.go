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
	command := "ls -la"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	log.Printf("Executing remote command: %s", command)
	fmt.Printf("📡 Connecting to device: %s\n", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")

	// Ejecutar comando remoto usando el prefijo "COMMAND:"
	// El cliente detectará automáticamente que es un comando y enviará type="command"
	rows, err := db.Query("COMMAND:" + command)
	if err != nil {
		log.Fatal("Error executing command:", err)
	}
	defer rows.Close()

	fmt.Println("\n⚡ --- Command Results ---")
	fmt.Printf("Command: %s\n", command)
	fmt.Println("----------------------------------------------------")

	// Procesar resultados - mostrar cada línea tal como viene
	lineCount := 0
	for rows.Next() {
		var output string
		if err := rows.Scan(&output); err != nil {
			log.Fatal("Error scanning result:", err)
		}
		fmt.Println(output)
		lineCount++
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating results:", err)
	}

	fmt.Printf("\n✅ Command executed successfully (%d lines output)\n", lineCount)
	fmt.Println("\n💡 Tips:")
	fmt.Println("   - You can specify a command as argument:")
	fmt.Println("     go run main.go \"ps aux\"")
	fmt.Println("   - The client automatically detects COMMAND: prefix")
	fmt.Println("   - Server will process this as type='command'")
	fmt.Println("   - Use quotes for commands with spaces")
	fmt.Println("   - Multi-line output is preserved")
}
