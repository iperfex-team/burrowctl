package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lordbasex/burrowctl/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar señales para cerrar gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Closing server...")
		cancel()
	}()

	// Configuración del pool de conexiones
	pool := &server.PoolConfig{
		MaxIdleConns:    5,
		MaxOpenConns:    15,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Crear el handler con credenciales hardcodeadas
	h := server.NewHandler(
		"fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb",     // Device ID
		"amqp://burrowuser:burrowpass123@localhost:5672/",                      // RabbitMQ URI
		"burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb?parseTime=true", // MariaDB DSN
		"open", // Modo de conexión: "open" para pool de conexiones
		pool,   // Configuración del pool
	)

	log.Println("Iniciando servidor burrowctl...")
	log.Println("Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	log.Println("RabbitMQ: localhost:5672")
	log.Println("MariaDB: localhost:3306/burrowdb")

	if err := h.Start(ctx); err != nil {
		log.Fatal("Error starting server:", err)
	}

	log.Println("Server closed")
}
