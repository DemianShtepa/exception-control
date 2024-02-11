package database

import (
	"context"
	"fmt"
	"github.com/DemianShtepa/exception-control/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Database struct {
	log  *slog.Logger
	Pool *pgxpool.Pool
}

func New(ctx context.Context, log *slog.Logger, cfg config.DBConfig) (*Database, error) {
	pgxCgx, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Connection,
			cfg.Port,
			cfg.Database,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCgx)
	if err != nil {
		return nil, fmt.Errorf("failed to init pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping pgx pool: %w", err)
	}

	return &Database{
		log:  log,
		Pool: pool,
	}, nil
}

func (db *Database) Stop() {
	db.log.Info("stopping database connection")

	db.Pool.Close()
}
