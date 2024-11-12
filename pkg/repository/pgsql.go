package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGRepository struct {
	mu sync.Mutex
	pool *pgxpool.Pool
}

// New создает новый экземпляр PGRepository, устанавливая соединение с базой данных
// по заданной строке подключения. Если соединение не удается установить, возвращает ошибку.
func New(connectionString string) (*PGRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), connectionString)

	if err != nil {
		return nil, err
	}

	return &PGRepository{mu: sync.Mutex{}, pool: pool}, nil
}