package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(connStr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать пул: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("База данных не отвечает на пинг: %w", err)
	}

	return dbPool, nil
}
