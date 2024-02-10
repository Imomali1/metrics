package server

import (
	"context"
	"errors"
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	store "github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/Imomali1/metrics/internal/tasks"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Run() {
	var cfg Config
	Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	storage, err := initStorage(cfg)
	if err != nil {
		log.Logger.Info().Err(err).Msg("failed to initialize storage")
		return
	}

	repo := repository.New(storage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:         log,
		ServiceManager: service,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}
	go func() {
		err = server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Logger.Info().Err(err).Msg("failed to listen and serve http server")
		}
	}()

	var wg sync.WaitGroup
	if cfg.FileStoragePath != "" && cfg.StoreInterval != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = tasks.WriteMetricsToFile(ctx, storage, cfg.FileStoragePath, cfg.StoreInterval); err != nil {
				log.Logger.Info().Err(err).Msg("error in writing metrics to file")
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	cancel()
	wg.Wait()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctxShutdown); err != nil {
		log.Logger.Info().Err(err).Msg("error in shutting down server")
	} else {
		log.Logger.Info().Msg("server stopped successfully")
	}
}

func initStorage(cfg Config) (*store.Storage, error) {
	if cfg.FileStoragePath == "" {
		return store.NewStorage()
	}

	var storageOptions []store.OptionsStorage

	if cfg.StoreInterval == 0 {
		storageOptions = append(storageOptions, store.WithFileStorage(cfg.FileStoragePath))
	}

	if cfg.Restore {
		storageOptions = append(storageOptions, store.RestoreFile(cfg.FileStoragePath))
	}

	return store.NewStorage(storageOptions...)
}
