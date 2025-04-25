package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/devdahcoder/golang-todo-api/util"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     time.Duration
	MaxRetries      int
	RetryInterval   int
	MaxConnLifetime time.Duration
}

func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.Database, dc.SSLMode)
}

func NewDatabaseConfig(env *util.EnvConfig) (*sql.DB, error) {
	dbCfg := &DatabaseConfig{
		Host: env.GetEnv("DB_HOST", "localhost"),
		Port: env.GetEnvAsInt("DB_PORT", 5432),
		User: env.GetEnv("DB_USER", "postgres"),
		Password: env.GetEnv("DB_PASSWORD", "postgres"),
		Database: env.GetEnv("DB_DATABASE", "myapp"),
		SSLMode: env.GetEnv("DB_SSLMODE", "disable"),
		
		MaxOpenConns:    env.GetEnvAsInt("MAX_OPEN_CONNS", 25),
		MaxIdleConns:    env.GetEnvAsInt("MAX_IDLE_CONNS", 25),
		MaxConnLifetime: time.Duration(env.GetEnvAsInt("MAX_CONN_LIFETIME_SECONDS", 1800)) * time.Second,
		MaxIdleTime:     time.Duration(env.GetEnvAsInt("MAX_IDLE_TIME_SECONDS", 1800)) * time.Second,
		MaxRetries:      env.GetEnvAsInt("MAX_RETRIES", 3),
		RetryInterval:   env.GetEnvAsInt("RETRY_INTERVAL_SECONDS", 5),
	}

	var db *sql.DB
	var err error

	for i := 0; i < dbCfg.MaxRetries; i++ {
		db, err = sql.Open("postgres", dbCfg.ConnectionString())
		if err != nil {
			log.Printf("Attempt %d: Error connecting to the database: %v", i+1, err)
			if i < dbCfg.MaxRetries-1 {
				time.Sleep(time.Duration(dbCfg.RetryInterval) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", dbCfg.MaxRetries, err)
		}
		break
	}

	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	db.SetConnMaxLifetime(dbCfg.MaxConnLifetime)
	db.SetConnMaxIdleTime(dbCfg.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging the database at %s:%d: %v", dbCfg.Host, dbCfg.Port, err)
	}

	log.Printf("Successfully connected to PostgreSQL database at %s:%d", dbCfg.Host, dbCfg.Port)

	return db, nil
}
