package service

import (
	"context"
	"fmt"

	"github.com/Bitummit/booking_api/internal/models"
)

type (
	HotelService struct {
		Storage HotelStorage
	}

	HotelStorage interface {
		CreateTag(ctx context.Context, tag models.Tag) (int64, error)
		ListTags(ctx context.Context) ([]models.Tag, error)
		DeleteTag(ctx context.Context, id int64) error
		CreateCity(ctx context.Context, city models.City) (int64, error)
		ListCities(ctx context.Context) ([]models.City, error)
		DeleteCity(ctx context.Context, id int64) error
		CreateHotel(ctx context.Context, hotel models.Hotel, cityName string, tags []string) (int64, error)
	}
)

func New(storage HotelStorage) *HotelService {
	return &HotelService{
		Storage: storage,
	}
}

func (s *HotelService) CreateTag(ctx context.Context, tag models.Tag) (int64, error) {
	id, err := s.Storage.CreateTag(ctx, tag)
	if err != nil {
		return 0, fmt.Errorf("creating new tag: %w", err)
	}
	return id, nil
}

func (s *HotelService) CreateCity(ctx context.Context, city models.City) (int64, error) {
	id, err := s.Storage.CreateCity(ctx, city)
	if err != nil {
		return 0, fmt.Errorf("creating new city: %w", err)
	}
	return id, nil
}

func (s *HotelService) ListTags(ctx context.Context) ([]models.Tag, error) {
	tags, err := s.Storage.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all tags: %w", err)
	}
	return tags, nil
}

func (s *HotelService) ListCities(ctx context.Context) ([]models.City, error) {
	cities, err := s.Storage.ListCities(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all cities: %w", err)
	}
	return cities, nil
}

func (s *HotelService) DeleteTag(ctx context.Context, id int64) error {
	err := s.Storage.DeleteTag(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting tag: %w", err)
	}
	return nil
}

func (s *HotelService) DeleteCity(ctx context.Context, id int64) error {
	err := s.Storage.DeleteCity(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting city: %w", err)
	}
	return nil
}

func (s *HotelService) CreateHotel(ctx context.Context, hotel models.Hotel, cityName string, tags []string) (int64, error) {
	// get city obj
	// city, err := s.Storage.getCity(ctx, cityName)
	// if err != nil {
	// 	return 0, fmt.Errorf("creating hotel: %w", err)
	// }
	// get tag objects
	//create city
	// hotel.CityId = city.Id
	hotelID, err := s.Storage.CreateHotel(ctx, hotel, cityName, tags)
	if err != nil {
		return hotelID, fmt.Errorf("creating hotel: %w", err)
	}
	return hotelID, nil
}