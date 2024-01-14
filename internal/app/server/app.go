package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"net/http"
)

func Run() error {
	if err := logger.InitLogger(); err != nil {
		return err
	}

	var cfg Config
	Parse(&cfg)

	memStorage := storage.NewStorage()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		ServiceManager: service,
	})

	logger.Log.Infow("Running server...", "address", cfg.ServerAddress)
	return http.ListenAndServe(cfg.ServerAddress, handler)
}
