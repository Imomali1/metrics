package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	stdLog "log"
	"net/http"
	"os"
)

func Run() {
	// Load configurations
	var cfg Config
	Parse(&cfg)

	// Create new logger instance
	log, err := logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		stdLog.Fatal(err)
	}

	handler := newHandler()
	log.
		Info().
		Str("address", cfg.ServerAddress).
		Msg("Running Server...")

	err = http.ListenAndServe(cfg.ServerAddress, handler)
	if err != nil {
		log.Fatal().Err(err)
	}
}

func newHandler() http.Handler {
	memStorage := storage.New()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		ServiceManager: service,
	})
	return handler
}
