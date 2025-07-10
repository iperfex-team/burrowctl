package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lordbasex/burrowctl/client"
	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// DSN con credenciales hardcodeadas para RabbitMQ
	dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true"

	// Verificar argumentos de línea de comandos
	command := "ls -la"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	fmt.Printf("📡 Connecting to device: %s\n", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	log.Printf("Executing remote command: %s", command)

	// Demonstrate both approaches
	fmt.Println("\n🔄 Demonstrating both approaches:")
	fmt.Println("1. Extended Client (Recommended)")
	fmt.Println("2. Standard database/sql Interface")
	fmt.Println("===========================================")

	// Approach 1: Extended Client (Recommended)
	fmt.Println("\n🚀 Method 1: Extended Client")
	demonstrateExtendedClient(dsn, command)

	// Approach 2: Standard database/sql (Legacy)
	fmt.Println("\n🔧 Method 2: Standard database/sql Interface")
	demonstrateStandardClient(dsn, command)

	fmt.Println("\n💡 Tips:")
	fmt.Println("   - Extended client provides cleaner interface")
	fmt.Println("   - Structured results with CommandResult type")
	fmt.Println("   - Better error handling and metadata")
	fmt.Println("   - Full compatibility with database/sql")
	fmt.Println("   - You can specify a command as argument:")
	fmt.Println("     go run main.go \"ps aux\"")
}

// demonstrateExtendedClient shows the new extended client approach
func demonstrateExtendedClient(dsn, command string) {
	// Create extended client
	bc, err := client.NewBurrowClient(dsn)
	if err != nil {
		log.Fatal("Error creating extended client:", err)
	}
	defer bc.Close()

	// Execute command using extended client
	result, err := bc.ExecCommand(command)
	if err != nil {
		log.Fatal("Error executing command:", err)
	}

	fmt.Printf("Command: %s\n", result.Command)
	fmt.Printf("Exit Code: %d\n", result.ExitCode)
	fmt.Printf("Executed At: %s\n", result.ExecutedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("----------------------------------------------------")

	// Display stdout
	if len(result.Stdout) > 0 {
		fmt.Println("📤 STDOUT:")
		for _, line := range result.Stdout {
			fmt.Println(line)
		}
	}

	// Display stderr
	if len(result.Stderr) > 0 {
		fmt.Println("📥 STDERR:")
		for _, line := range result.Stderr {
			fmt.Println(line)
		}
	}

	fmt.Printf("✅ Command executed successfully (%d stdout lines, %d stderr lines)\n", 
		len(result.Stdout), len(result.Stderr))
}

// demonstrateStandardClient shows the traditional database/sql approach
func demonstrateStandardClient(dsn, command string) {
	// Abrir conexión usando el driver rabbitsql
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer db.Close()

	// Ejecutar comando remoto usando el prefijo "COMMAND:"
	// El cliente detectará automáticamente que es un comando y enviará type="command"
	rows, err := db.Query("COMMAND:" + command)
	if err != nil {
		log.Fatal("Error executing command:", err)
	}
	defer rows.Close()

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

	fmt.Printf("✅ Command executed successfully (%d lines output)\n", lineCount)
}
