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
	"github.com/lib/pq"
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

	err := s.DB.QueryRow(ctx, checkStmt, args).Scan(&id)
	if err == nil {
		return 0, fmt.Errorf("database error: %w", ErrorExists)
	}

	createStmt := CreateTagStmt
	err = s.DB.QueryRow(ctx, createStmt, args).Scan(&id)
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

	err := s.DB.QueryRow(ctx, checkStmt, args).Scan(&id)
	if err == nil {
		return 0, fmt.Errorf("database error: %w", ErrorExists)
	}

	stmt := CreateCityStmt
	err = s.DB.QueryRow(ctx, stmt, args).Scan(&id)
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

func (s *Storage) DeleteTag(ctx context.Context, id int64) error {
	stmt := DeleteTagStmt
	args := pgx.NamedArgs{
		"id": id,
	}

	resp, err := s.DB.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("no such tag")
	}
	if resp.RowsAffected() == 0 {
		return fmt.Errorf("deleting: %w", ErrorNotExists)
	}

	return nil
}

func (s *Storage) DeleteCity(ctx context.Context, id int64) error {
	stmt := DeleteCityStmt
	args := pgx.NamedArgs{
		"id": id,
	}

	resp, err := s.DB.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("deleting err: %w", err)
	}
	if resp.RowsAffected() == 0 {
		return fmt.Errorf("deleting: %w", ErrorNotExists)
	}
	
	return nil
}

func (s *Storage) CreateHotel(ctx context.Context, hotel models.Hotel, cityName string, tagNames []string) (int64, error) {
	var id int64
	// dont work!!!!!
	// stmt := CheckHotelNameUniqueStmt
	// args := pgx.NamedArgs{
	// 	"name": hotel.Name,
	// }
	// err := s.DB.QueryRow(ctx, stmt, args).Scan(&id)
	// if err != nil {
	// 	return 0, fmt.Errorf("database error: %w", ErrorExists)
	// }

	stmt := CreateHotelStmt
	args := pgx.NamedArgs{
		"name": hotel.Name,
		"desc": hotel.Desc,
		"city_name": cityName,
	}
	err := s.DB.QueryRow(ctx, stmt, args).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("inserting error: %w", ErrorInsertion)
		}
		return 0, fmt.Errorf("database internal error: %w", err)
	}

	tagsID, err := s.GetTagsID(ctx, tagNames)
	if err != nil {
		return 0, fmt.Errorf("getting tag ids: %w", err)
	}
	// create hotel_tags
	for _,tagID := range tagsID {
		err = s.CreateTagHotel(ctx, id, tagID)
		if err != nil {
			return 0, fmt.Errorf("%w", err)
		}
	}
	
	return id, nil
}

func (s *Storage) GetTagsID(ctx context.Context, tags []string) ([]int64, error) {
	var tagsID []int64
	stmt := GetMultipleTagsStmt
	args := pgx.NamedArgs{
		"tag_array": pq.Array(tags),
	}
	
	rows, err := s.DB.Query(ctx, stmt, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("query error: %w", ErrorNotExists)
		}
		return nil, fmt.Errorf("database internal error: %w", err)
	}
	for rows.Next() {
		var tagID int64
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		tagsID = append(tagsID, tagID)
	}
	if len(tagsID) == 0 {
		return nil, fmt.Errorf("tag not exist:%w", ErrorTagNotExists)
	}
	return tagsID, nil
}

func (s *Storage) CreateTagHotel(ctx context.Context, hotelID int64, tagID int64) error {
	stmt := CreateTagHotelStmt
	args := pgx.NamedArgs{
		"tag_id": tagID,
		"hotel_id": hotelID,
	}
	resp, err := s.DB.Exec(ctx, stmt, args)
	if err != nil {
		return fmt.Errorf("isnerting hotel_tag: %w", err)
	}

	if resp.RowsAffected() == 0 {
		return fmt.Errorf("not inserted: %w", ErrorInsertion)
	}

	return nil
}