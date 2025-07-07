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

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	pool := &server.PoolConfig{
		MaxIdleConns:    5,
		MaxOpenConns:    15,
		ConnMaxLifetime: 5 * time.Minute,
	}

	h := server.NewHandler(
		"fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb",
		"amqp://burrowuser:burrowpass123@localhost:5672/",
		"burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb?parseTime=true",
		"open", // "" defaults to open
		pool,
	)

	if err := h.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
