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

	// Configure signals to close gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Closing server...")
		cancel()
	}()

	// Configuration of the connection pool
	pool := &server.PoolConfig{
		MaxIdleConns:    5,
		MaxOpenConns:    15,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Create the handler with hardcoded credentials
	h := server.NewHandler(
		"fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb",   // Device ID
		"amqp://burrowuser:burrowpass123@rabbitmq:5672/",                     // RabbitMQ URI
		"burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb?parseTime=true", // MariaDB DSN
		"open", // Connection mode: "open" for connection pool
		pool,   // Configuration of the pool
	)

	log.Println("Starting burrowctl server...")
	log.Println("Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	log.Println("RabbitMQ: rabbitmq:5672")
	log.Println("MariaDB: mariadb:3306/burrowdb")

	if err := h.Start(ctx); err != nil {
		log.Fatal("Error starting server:", err)
	}

	log.Println("Server closed")
}
