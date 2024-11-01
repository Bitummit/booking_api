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
	var id int64
	checkStmt := CheckTagNameUniqueStmt
	args := pgx.NamedArgs{
		"name": tag.Name,
	}

	row := s.DB.QueryRow(ctx, checkStmt, args).Scan(&id)
	if row == nil {
		return 0, fmt.Errorf("database error: %w", ErrorExists)
	}

	createStmt := CreateTagStmt
	err := s.DB.QueryRow(ctx, createStmt, args).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("database error: %w", ErrorInsertion)
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

func (s *Storage) CreateCity(ctx context.Context, city models.City) (int64, error) {
	var id int64
	checkStmt := CheckCityNameUniqueStmt
	args := pgx.NamedArgs{
		"name": city.Name,
	}

	row := s.DB.QueryRow(ctx, checkStmt, args).Scan(&id)
	if row == nil {
		return 0, fmt.Errorf("database error: %w", ErrorExists)
	}

	stmt := CreateCityStmt
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

func (s *Storage) ListTags(ctx context.Context) ([]models.Tag, error) {
	stmt := ListTagsStmt
	var tags []models.Tag

	rows, err := s.DB.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("fetching data: %w", err)
	}

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("fetching data: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) ListCities(ctx context.Context) ([]models.City, error) {
	stmt := ListCitiesStmt
	var cities []models.City

	rows, err := s.DB.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("fetching data: %w", err)
	}

	for rows.Next() {
		var city models.City
		err = rows.Scan(&city.Id, &city.Name)
		if err != nil {
			return nil, fmt.Errorf("fetching data: %w", err)
		}

		cities = append(cities, city)
	}

	return cities, nil
}
