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
		CreateCity(ctx context.Context, city models.City) (int64, error)
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