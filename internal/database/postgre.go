package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(dbURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	connURL := dbURL

	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to DB")
	return pool, nil
}
