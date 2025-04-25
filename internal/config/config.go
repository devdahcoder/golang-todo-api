package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/devdahcoder/golang-todo-api/internal/database"
	"github.com/devdahcoder/golang-todo-api/util"
)

type Config struct {
    ServerAddress   string
    JWTSecret       string
    JWTExpiryHours  int
    Environment     string
    LogLevel        string
    CacheEnabled    bool
    ShutdownTimeout time.Duration
	Db *sql.DB
}

func Load() (*Config, error) {
	env, err := util.NewEnvConfig()
    if err != nil {
        log.Printf("Warning: Error loading environment variables: %v", err)
        env = &util.EnvConfig{}
    }

	db, err := database.NewDatabaseConfig(env)
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
    }
    
    return config, nil
}

func (c *Config) Close() error {
    if c.Db != nil {
        return c.Db.Close()
    }
    return nil
}