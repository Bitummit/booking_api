package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Bitummit/booking_api/pkg/config"
	"github.com/go-chi/chi/v5"
)

type (
	HTTPServer struct {
		Cfg *config.Config
		Log *slog.Logger
		Storage Storage
		Router chi.Router
	}

	Storage interface {

	}
)

func New(cfg *config.Config, log *slog.Logger, storage Storage) *HTTPServer{
	router := chi.NewRouter()

	return &HTTPServer{
		Cfg: cfg,
		Log: log,
		Storage: storage,
		Router: router,
	}
}

func (s *HTTPServer) Run(ctx context.Context, wg *sync.WaitGroup) error {
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