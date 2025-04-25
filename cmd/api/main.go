package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devdahcoder/golang-todo-api/internal/config"
	"github.com/devdahcoder/golang-todo-api/internal/server"
)

func main() {
	cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    defer cfg.Close()
	
	app := server.New(cfg)
    
    go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
        if err := app.Listen(cfg.ServerAddress); err != nil {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    sig := <-quit
    
    log.Printf("Received shutdown signal: %v", sig)

    _, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    log.Println("Initiating graceful shutdown...")
    if err := app.Shutdown(); err != nil {
        log.Printf("Warning: Graceful shutdown failed: %v", err)
    }
    
    log.Println("Server shutdown complete")

}