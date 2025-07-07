package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// DSN con credenciales hardcodeadas para RabbitMQ
	dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s"

	// Abrir conexión usando el driver rabbitsql
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer db.Close()

	log.Println("Ejecutando consulta SELECT id, name FROM users...")

	// Ejecutar query
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal("Error en la consulta:", err)
	}
	defer rows.Close()

	fmt.Println("\n--- Resultados ---")
	fmt.Printf("%-5s %-30s\n", "ID", "Nombre")
	fmt.Println("------------------------------------")

	// Procesar resultados
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Error al escanear resultado:", err)
		}
		fmt.Printf("%-5d %-30s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error durante la iteración:", err)
	}

	fmt.Println("\n✅ Consulta completada exitosamente")
}
