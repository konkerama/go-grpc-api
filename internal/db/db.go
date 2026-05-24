package db

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DB_URL = "postgres://postgres:postgres@localhost:15432/postgres"
)

type DBConfig struct {
	PGPool *pgxpool.Pool
}

func NewDBConfig() *DBConfig {
	var dbConfig DBConfig
	dbConfig.PGPool = pgInit()
	return &dbConfig
}

func (r *DBConfig) CloseDB() (bool, error) {
	log.Printf("Closing postgres connection")
	if r.PGPool != nil {
		r.PGPool.Close()
	}
	return true, nil
}

// To make it accessible from outside you have to capitalize it
func pgInit() *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(DB_URL)
	if err != nil {
		slog.Error("Unable to parse database URL", "error", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = time.Minute * 5

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		slog.Error("Postgres: Unable to create connection pool", "error", err)
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	retryTicker := time.NewTicker(2 * time.Second)
	defer retryTicker.Stop()

	for {
		err = pool.Ping(context.Background())
		if err == nil {
			return pool
		}

		select {
		case <-timeoutCtx.Done():
			slog.Error("Postgres: Timed out waiting for database connection after 1 minute", "error", err)
		case <-retryTicker.C:
			slog.Warn("Postgres: Database not ready yet, retrying", "error", err)
		}
	}
}
