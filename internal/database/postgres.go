package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "postgres")
	dbName := config.GetEnv("DB_NAME", "rms")
	sslmode := config.GetEnv("DB_SSLMODE", "disable")
	maxConnsStr := config.GetEnv("DB_MAX_CONNS", "10")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslmode,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	maxConns, err := strconv.ParseInt(maxConnsStr, 10, 32)
	if err != nil {
		maxConns = 10
	}
	cfg.MaxConns = int32(maxConns)
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
