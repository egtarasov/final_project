package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	connString = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
)

func NewDb(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return new(pool), nil
}
