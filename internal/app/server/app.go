package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Imomali1/metrics/internal/pkg/file"

	_ "net/http/pprof"

	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/tasks"
	"github.com/Imomali1/metrics/internal/usecase"
)

const (
	_timeout          = 1 * time.Second
	_htmlTemplatePath = "static/templates/*.html"
)

func Run(cfg Config, log logger.Logger) {
	store, err := storage.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize storage")
	}

	if cfg.Restore {
		err = file.RestoreMetrics(context.Background(), cfg.FileStoragePath, store)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to restore metrics")
		}
	}

	var syncFileWriter file.SyncFileWriter
	if cfg.StoreInterval == 0 {
		syncFileWriter, err = file.NewSyncMetricsWriter(cfg.FileStoragePath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize sync file writer")
		}
	}

	repo := repository.New(store, syncFileWriter)
	uc := usecase.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:           log,
		UseCase:          uc,
		Cfg:              cfg.API,
		HTMLTemplatePath: _htmlTemplatePath,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to listen and serve http server")
		}
	}()

	var wg sync.WaitGroup
	if cfg.FileStoragePath != "" && cfg.StoreInterval != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = tasks.WriteMetricsToFile(ctx, store, cfg.FileStoragePath, cfg.StoreInterval); err != nil {
				log.Error().Err(err).Msg("error in writing metrics to file")
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	cancel()
	wg.Wait()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), _timeout)
	defer cancel()

	if err = server.Shutdown(ctxShutdown); err != nil {
		log.Fatal().Err(err).Msg("error in shutting down server")
	}

	log.Info().Msg("server stopped successfully")
}
