package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrorInsertion = errors.New("can not insert")
var ErrorExists = errors.New("already exists")


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
	stmt := CreateTagStmt
	args := pgx.NamedArgs{
		"name": tag.Name,
	}
	var id int64

	err := s.DB.QueryRow(ctx, stmt, args).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("database error: %w", ErrorInsertion)
		} else if errors.Is(err, pgx.ErrTooManyRows) {
			return 0, fmt.Errorf("database error: %w", ErrorExists)
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

func (s *Storage) CreateCity(ctx context.Context, city models.City) (int64, error) {
	stmt := CreateTagStmt
	args := pgx.NamedArgs{
		"name": city.Name,
	}
	var id int64

	err := s.DB.QueryRow(ctx, stmt, args).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("database error: %w", ErrorInsertion)
		} else if errors.Is(err, pgx.ErrTooManyRows) {
			return 0, fmt.Errorf("database error: %w", ErrorExists)
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}