package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database URL: %w", err)
    }

    // Connection pool settings - adjust based on your workload
    config.MaxConns = 25                      // Maximum connections in the pool
    config.MinConns = 5                       // Keep at least 5 connections ready
    config.MaxConnLifetime = time.Hour        // Recycle connections hourly
    config.MaxConnIdleTime = 30 * time.Minute // Close idle connections after 30 min
    config.HealthCheckPeriod = time.Minute    // Check connection health every minute

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    // Verify the connection actually works
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return pool, nil
}