package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Bitummit/booking_api/internal/middlewares"
	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/service"
	"github.com/Bitummit/booking_api/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	HTTPServer struct {
		Cfg *config.Config
		Log *slog.Logger
		HotelService HotelService
		Router chi.Router
	}

	HotelService interface {
		CreateTag(ctx context.Context, tag models.Tag) (int64, error)
		ListTags(ctx context.Context) ([]models.Tag, error)
		DeleteTag(ctx context.Context, id int64) error
		CreateCity(ctx context.Context, city models.City) (int64, error)
		ListCities(ctx context.Context) ([]models.City, error)
		DeleteCity(ctx context.Context, id int64) error
	}
)

func New(cfg *config.Config, log *slog.Logger, storage service.HotelStorage) *HTTPServer{
	router := chi.NewRouter()
	hotelService := service.New(storage)
	return &HTTPServer{
		Cfg: cfg,
		Log: log,
		HotelService: hotelService,
		Router: router,
	}
}

func (s *HTTPServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(middlewares.SetJSONContentType)
	
	s.Router.Post("/tag", s.CreateTagHandler)
	s.Router.Get("/tag", s.ListTagsHandler)
	s.Router.Delete("/tag/{id}", s.DeleteTagHandler)
	s.Router.Post("/city", s.CreateCityHandler)
	s.Router.Get("/city", s.ListCityHandler)
	s.Router.Delete("/city/{id}", s.DeleteCityHandler)

	errCh := make(chan error, 1)
	httpServer := &http.Server{
		Addr: s.Cfg.Address,
		Handler: s.Router,
		ReadTimeout: s.Cfg.Timeout,
		WriteTimeout: s.Cfg.Timeout,
		IdleTimeout: s.Cfg.IdleTimeout,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
			return
		}
	}()
	select {
	case err := <-errCh:
		return fmt.Errorf("listenning and serving: %w", err)
	case <-ctx.Done():
		s.Log.Info("Shutting down")
		httpServer.Shutdown(ctx)
		s.Log.Info("Server stopped")
		wg.Done()
	}
	return nil
	
}


// User:
// 	List hotels (add filters)
// 	Get hotel -> show list room_categories
// 	Create booking (auth)

// Admin:
// 	Create hotel
// 	Update user role (give role manager)
//	List, Create(done), delete tags -> 01.11.2024
// 	List, Create(done), delete city -> 01.11.2024

// Manager:
//	List own hotels, get hotel admin
//	Create, update, delete categories
//	Create, delete room
// 	Update hotel
//	