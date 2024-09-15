package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Imomali1/metrics/internal/pkg/file"

	_ "net/http/pprof"

	"crypto/rsa"

	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/cipher"
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

func Run(cfg Config, log logger.Logger) error {
	store, err := storage.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if cfg.Restore {
		err = file.RestoreMetrics(context.Background(), cfg.FileStoragePath, store)
		if err != nil {
			return fmt.Errorf("failed to restore metrics: %w", err)
		}
	}

	var syncFileWriter file.SyncFileWriter
	if cfg.StoreInterval == 0 {
		syncFileWriter, err = file.NewSyncMetricsWriter(cfg.FileStoragePath)
		if err != nil {
			return fmt.Errorf("failed to initialize sync file writer: %w", err)
		}
	}

	var privateKey *rsa.PrivateKey
	if cfg.PrivateKeyPath != "" {
		privateKey, err = cipher.UploadRSAPrivateKey(cfg.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to upload rsa private key: %w", err)
		}
	}

	repo := repository.New(store, syncFileWriter)
	uc := usecase.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:           log,
		UseCase:          uc,
		Cfg:              cfg.API,
		HTMLTemplatePath: _htmlTemplatePath,
		PrivateKey:       privateKey,
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

			err = tasks.WriteMetricsToFile(ctx,
				store,
				cfg.FileStoragePath,
				time.Duration(cfg.StoreInterval)*time.Second,
			)
			if err != nil {
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
		return fmt.Errorf("error in shutting down server: %w", err)
	}

	log.Info().Msg("server stopped successfully")

	return nil
}
