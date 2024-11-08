package run

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Bitummit/booking_api/internal/api/rest"
	"github.com/Bitummit/booking_api/internal/storage/postgresql"
	"github.com/Bitummit/booking_api/pkg/config"
	"github.com/Bitummit/booking_api/pkg/logger"
)


func Run() {
	wg := &sync.WaitGroup{}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()
	log := logger.NewLogger()
	log.Info("Config and logger inited")

	log.Info("Connecting database")
	storage, err := postgresql.New(ctx)
	if err != nil {
		log.Error("DB connection %v", logger.Err(err))
		return
	}
	log.Info("Database connected")

	wg.Add(1)
	log.Info("Starting http server")
	server, err := rest.New(cfg, log, storage)
	if err != nil {
		log.Error("starting server: ", logger.Err(err))
		storage.DB.Close()
		return
	}
	server.Start(ctx, wg)

	<-ctx.Done()
	wg.Wait()
	storage.DB.Close()
}
