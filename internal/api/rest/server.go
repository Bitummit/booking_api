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
	authclient "github.com/Bitummit/booking_api/internal/service/authClient"
	"github.com/Bitummit/booking_api/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	HTTPServer struct {
		Cfg *config.Config
		Log *slog.Logger
		HotelService HotelService
		AuthService *authclient.Client
		Router chi.Router
	}

	HotelService interface {
		CreateTag(ctx context.Context, tag models.Tag) (int64, error)
		ListTags(ctx context.Context) ([]models.Tag, error)
		DeleteTag(ctx context.Context, id int64) error
		CreateCity(ctx context.Context, city models.City) (int64, error)
		ListCities(ctx context.Context) ([]models.City, error)
		DeleteCity(ctx context.Context, id int64) error
		CreateHotel(ctx context.Context, hotel models.Hotel, cityName string, tags []string) (int64, error)
	}
)

func New(cfg *config.Config, log *slog.Logger, storage service.HotelStorage) (*HTTPServer, error){
	router := chi.NewRouter()
	hotelService := service.New(storage)

	auth, err := authclient.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &HTTPServer{
		Cfg: cfg,
		Log: log,
		HotelService: hotelService,
		AuthService: auth,
		Router: router,
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(middlewares.SetJSONContentType)

	s.Router.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.IsAdmin)

		r.Route("/tag", func(r chi.Router) {
			r.Post("/", s.CreateTagHandler)
			r.Get("/", s.ListTagsHandler)
			r.Delete("/{id}", s.DeleteTagHandler)
		})
		r.Route("/city", func(r chi.Router) {
			r.Post("/", s.CreateCityHandler)
			r.Get("/", s.ListCityHandler)
			r.Delete("/{id}", s.DeleteCityHandler)
		})
	})
	s.Router.Post("/hotel", s.CreateHotelHandler) // manager
	s.Router.Post("/signup", s.RegistrationHandler)
	s.Router.Post("/login", s.LoginHandler)

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
// 	Update user role (give role manager) (auth_service)
//	List, Create(done), delete tags -> 01.11.2024
// 	List, Create(done), delete city -> 01.11.2024

// Manager:
//	List own hotels, get hotel admin
//	Create, update, delete categories
//	Create, delete room
// 	Update hotel
//	Create hotel