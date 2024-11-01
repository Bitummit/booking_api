package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/service"
	"github.com/Bitummit/booking_api/internal/storage/postgresql"
	"github.com/Bitummit/booking_api/pkg/config"
	"github.com/Bitummit/booking_api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
		CreateCity(ctx context.Context, city models.City) (int64, error)
		ListCities(ctx context.Context) ([]models.City, error)
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
	// register middllewares
	s.Router.Post("/tag", s.CreateTagHandler)
	s.Router.Get("/tag", s.ListTagsHandler)
	s.Router.Post("/city", s.CreateCityHandler)
	s.Router.Get("/city", s.ListCityHandler)

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

// Admin role
func (s *HTTPServer) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTagRequest
	err := render.DecodeJSON(r.Body, &req)
	r.Body.Close()
	if err != nil {
		s.Log.Error("decoding request: %v", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse("bad request"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		err = err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	tag := models.Tag{
		Name: req.Name,
	}
	id, err := s.HotelService.CreateTag(r.Context(), tag)
	if err != nil {
		s.Log.Error("%v", logger.Err(err))
		if errors.Is(err, postgresql.ErrorInsertion){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("insertion error"))
			return
		} else if errors.Is(err, postgresql.ErrorExists){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("tag aready exists"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	s.Log.Info("New tag", slog.Int64("id", int64(id)))
	res := CreateTagResponse{
		Id: id,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}

// Admin role
func (s *HTTPServer) CreateCityHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCityRequest
	err := render.DecodeJSON(r.Body, &req)
	r.Body.Close()
	if err != nil {
		s.Log.Error("decoding request: %v", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse("bad request"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		err = err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	city := models.City{
		Name: req.Name,
	}
	id, err := s.HotelService.CreateCity(r.Context(), city)
	if err != nil {
		s.Log.Error("%v", logger.Err(err))
		if errors.Is(err, postgresql.ErrorInsertion){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("insertion error"))
			return
		} else if errors.Is(err, postgresql.ErrorExists){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("tag aready exists"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	s.Log.Info("New city", slog.Int64("id", int64(id)))
	res := CreateCityResponse{
		Id: id,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}

func (s *HTTPServer) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := s.HotelService.ListTags(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, ListTagResponse{
		Tags: tags,
	})
}

func (s *HTTPServer) ListCityHandler(w http.ResponseWriter, r *http.Request) {
	cities, err := s.HotelService.ListCities(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, ListCityResponse{
		Cities: cities,
	})
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