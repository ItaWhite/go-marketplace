package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(user, password, host, database string) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, host, database)
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	return db, nil
}
