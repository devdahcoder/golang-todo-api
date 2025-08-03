package util

import (
	"os"
	"strconv"

	"github.com/devdahcoder/golang-todo-api/pkg/logger"
	"github.com/joho/godotenv"
)

type EnvConfig struct{}

func NewEnvConfig(zapLogger *logger.Logger) (*EnvConfig, error) {
    err := godotenv.Load()
    if err != nil {
        if os.IsNotExist(err) {
            zapLogger.Warn("Environment variables file not found")
            return &EnvConfig{}, nil
        }
        return nil, err
    }
    return &EnvConfig{}, nil
}

func (c *EnvConfig) GetEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func (c *EnvConfig) GetEnvAsInt(key string, defaultValue int) int {
    if valueStr, exists := os.LookupEnv(key); exists {
        if value, err := strconv.Atoi(valueStr); err == nil {
            return value
        }
    }
    return defaultValue
}

func (c *EnvConfig) GetEnvAsBool(key string, defaultValue bool) bool {
    if valueStr, exists := os.LookupEnv(key); exists {
        if value, err := strconv.ParseBool(valueStr); err == nil {
            return value
        }
    }
    return defaultValue
}