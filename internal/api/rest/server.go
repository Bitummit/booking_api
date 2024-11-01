package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/service"
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
	// register endpoints
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
func (s *HTTPServer) CreateTag(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	s.Log.Info("New tag", slog.Int64("id", int64(id)))
	res := CreateTagResponse{
		Response: Response{
			Status: "OK",
		},
		Id: id,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}

// User:
// 	List hotels (add filters)
// 	Get hotel -> show list room_categories
// 	Create booking (auth)

// Admin:
// 	Create hotel
// 	Update user role (make manager)
//	List, Create, delete tags
// 	List, Create, delete city

// Manager:
//	List own hotels, get hotel admin
//	Create, update, delete categories
//	Create, delete room
// 	Update hotel
//	