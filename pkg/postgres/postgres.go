package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(connUrl string, maxConns int) (*pgxpool.Pool, error) {
	const op = "db.NewPostgresPool"

	conf, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conf.MaxConns = int32(maxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pool, nil
}
