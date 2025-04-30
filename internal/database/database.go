package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/devdahcoder/golang-todo-api/pkg/logger"
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
	zapLogger *logger.Logger
}

func (dc *DatabaseConfig) ConnectionString() string {
		dc.zapLogger.Info("Database connection string", zap.String("connection_string", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			dc.User, dc.Password, dc.Host, dc.Port, dc.Database, dc.SSLMode)))

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Database, dc.SSLMode)
}

func NewPqDatabaseConfig(env *util.EnvConfig, zapLogger *logger.Logger) (*sql.DB, error) {
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
		zapLogger: zapLogger,
	}

	var db *sql.DB
	var err error

	for i := 0; i < dbCfg.MaxRetries; i++ {
		db, err = sql.Open("postgres", dbCfg.ConnectionString())
		
		if err != nil {
			dbCfg.zapLogger.Error("Attempt to connect to the database failed", zap.Int("attempt", i+1), zap.Error(err))
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

	dbCfg.zapLogger.Info("Database connection established", zap.String("host", dbCfg.Host), zap.Int("port", dbCfg.Port))

	return db, nil
}

func (dbCfg *DatabaseConfig) NewPgxDatabaseConfig(env *util.EnvConfig, zapLogger *logger.Logger) (*pgxpool.Pool, error) {
	dbCfg = &DatabaseConfig{
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
		zapLogger: zapLogger,
	}

	var pool *pgxpool.Pool
	var err error

	config, err := pgxpool.ParseConfig(dbCfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %v", err)
	}

	config.MaxConns = int32(dbCfg.MaxOpenConns)
	config.MaxConnLifetime = dbCfg.MaxConnLifetime
	config.MaxConnIdleTime = dbCfg.MaxIdleTime

	for i := 0; i < dbCfg.MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.NewWithConfig(ctx, config)
		cancel()

		if err != nil {
			dbCfg.zapLogger.Error("Attempt to connect to the database failed", 
				zap.Int("attempt", i+1), 
				zap.Error(err))
			
			if i < dbCfg.MaxRetries-1 {
				time.Sleep(time.Duration(dbCfg.RetryInterval) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", dbCfg.MaxRetries, err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		err = pool.Ping(ctx)
		cancel()

		if err != nil {
			pool.Close()
			dbCfg.zapLogger.Error("Failed to ping database", 
				zap.Int("attempt", i+1), 
				zap.Error(err))
			
			if i < dbCfg.MaxRetries-1 {
				time.Sleep(time.Duration(dbCfg.RetryInterval) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to ping database after %d attempts: %v", dbCfg.MaxRetries, err)
		}

		break
	}

	dbCfg.zapLogger.Info("Database connection established with pgx", 
		zap.String("host", dbCfg.Host), 
		zap.Int("port", dbCfg.Port))

	return pool, nil
}