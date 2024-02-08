package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	store "github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"net/http"
	"os"
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

	if cfg.FileStoragePath != "" && cfg.StoreInterval != 0 {
		done := make(chan struct{}, 1)
		defer func() {
			done <- struct{}{}
			close(done)
		}()
		go storeMetricsToFilePeriodically(log, storage, cfg.StoreInterval, done)
	}

	err = http.ListenAndServe(cfg.ServerAddress, handler)
	if err != nil {
		log.Logger.Info().Err(err).Msg("failed to listen and serve http server")
	}
}

func initStorage(cfg Config) (*store.Storage, error) {
	if cfg.FileStoragePath == "" {
		return store.NewStorage()
	}
	var storageOptions []store.OptionsStorage
	if cfg.StoreInterval == 0 {
		storageOptions = append(storageOptions,
			store.WithFileStorage(cfg.FileStoragePath),
			store.SyncWriteFile())
	}

	if cfg.Restore {
		storageOptions = append(storageOptions, store.RestoreFile(cfg.FileStoragePath))
	}

	return store.NewStorage(storageOptions...)
}

func storeMetricsToFilePeriodically(log logger.Logger, storage *store.Storage, interval int, done chan struct{}) {
	storeTicker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-storeTicker.C:
			err := storage.File.WriteAllMetrics()
			if err != nil {
				log.Logger.Info().Err(err).Msg("cannot write metrics to file")
			}
		case <-done:
			return
		}
	}
}
