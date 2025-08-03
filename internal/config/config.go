package config

import (
	// "database/sql"
	"fmt"
	"os"
	"time"

	"github.com/devdahcoder/golang-todo-api/internal/database"
	"github.com/devdahcoder/golang-todo-api/pkg/logger"
	"github.com/devdahcoder/golang-todo-api/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Config struct {
    ServerAddress   string
    JWTSecret       string
    JWTExpiryHours  int
    Environment     string
    LogLevel        string
    CacheEnabled    bool
    ShutdownTimeout time.Duration
	Db *pgxpool.Pool
    ZapLogger *logger.Logger
}

func Load() (*Config, error) {
    zapLogger := logger.NewLogger(logger.LoggerConfig{
		LogLevel: "debug",
		Output:   os.Stdout,
		SkipPaths: []string{"/health"},
	})
	defer zapLogger.Close()
    
	env, err := util.NewEnvConfig(zapLogger)
    if err != nil {
        zapLogger.Warn("Error loading environment variables", zap.Error(err))
        env = &util.EnvConfig{}
    }

	db, err := database.NewPgxDatabaseConfig(env, zapLogger)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }
	
    config := &Config{
        ServerAddress:   env.GetEnv("SERVER_ADDRESS", ":8080"),
        JWTSecret:       env.GetEnv("JWT_SECRET", "secret-key"),
        JWTExpiryHours:  env.GetEnvAsInt("JWT_EXPIRY_HOURS", 24),
        Environment:     env.GetEnv("ENVIRONMENT", "development"),
        LogLevel:        env.GetEnv("LOG_LEVEL", "info"),
        CacheEnabled:    env.GetEnvAsBool("CACHE_ENABLED", false),
        ShutdownTimeout: time.Duration(env.GetEnvAsInt("SHUTDOWN_TIMEOUT_SECONDS", 5)) * time.Second,
		Db: db,
        ZapLogger: zapLogger,
    }
    
    return config, nil
}

func (c *Config) Close() {
    if c.Db != nil {
        c.Db.Close()
    }
}