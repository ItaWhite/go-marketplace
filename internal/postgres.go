package internal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDb(url string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
