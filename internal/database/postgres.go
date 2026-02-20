package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	poolCfg.MaxConns = 15
	poolCfg.MinConns = 3
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
