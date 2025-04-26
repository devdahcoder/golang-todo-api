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
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.Load()
    if err != nil {
		cfg.ZapLogger.Fatal("Failed to load config", zap.Error(err))
    }
    defer cfg.Close()
	
	app := server.NewServer(cfg)
    
    go func() {
		cfg.ZapLogger.Info("Starting server", zap.String("address", cfg.ServerAddress))
        if err := app.Listen(":" + cfg.ServerAddress); err != nil {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    sig := <-quit
    
	cfg.ZapLogger.Info("Received shutdown signal", zap.String("signal", sig.String()))

    _, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

	cfg.ZapLogger.Info("Shutting down server...")
    if err := app.Shutdown(); err != nil {
		cfg.ZapLogger.Error("Error shutting down server", zap.Error(err))
    }
    
	cfg.ZapLogger.Info("Server shutdown complete")

}