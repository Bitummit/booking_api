package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bitummit/booking_api/internal/models"
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

func (s *Storage) CreateTag(ctx context.Context, tag models.Tag) (int64, error) {
	return 1, nil
}

func (s *Storage) CreateCity(ctx context.Context, city models.City) (int64, error) {
	return 1, nil
}