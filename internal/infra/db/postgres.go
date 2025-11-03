package db

import (
	"context"
	"ecommerce-app/configs"
	"ecommerce-app/internal/pkg/database/sqlc"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	// Get the database URL from the config
	dbURL := getDatabaseURL(cfg)

	// Check if the database URL is empty
	if dbURL == "" {
		return nil, fmt.Errorf("database URL is not set")
	}

	// Connect to the PostgreSQL database
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Test the connection by pinging the database
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %v", err)
	}

	// Return the pool if everything is fine
	return pool, nil
}

func NewQueries(pool *pgxpool.Pool) *sqlc.Queries {
	return sqlc.New(pool)
}

func getDatabaseURL(cfg *configs.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}
