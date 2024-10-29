package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context) (*Storage, error){
	dbURL := os.Getenv("DB_URL")
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}
	return &Storage{DB: db}, nil
}

