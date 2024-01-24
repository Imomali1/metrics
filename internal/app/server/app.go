package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"net/http"
	"os"
)

func Run() {
	var cfg Config
	Parse(&cfg)
	l := logger.NewLogger(os.Stdout, cfg.LogLevel, "server")
	memStorage := storage.New()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:         l,
		ServiceManager: service,
	})

	l.Logger.Info().
		Str("address", cfg.ServerAddress).
		Msg("Running Server...")

	err := http.ListenAndServe(cfg.ServerAddress, handler)
	if err != nil {
		l.Logger.Fatal().
			Err(err).
			Msg("Server is stopped.")
	}
}
