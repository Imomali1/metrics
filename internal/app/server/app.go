package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/Imomali1/metrics/internal/tasks/file_storage"
	"net/http"
	"os"
)

func Run() {
	var cfg Config
	Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	memStorage := storage.NewStorage()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:         log,
		ServiceManager: service,
	})

	go file_storage.RunTask(memStorage.MetricStorage)

	err := http.ListenAndServe(cfg.ServerAddress, handler)
	if err != nil {
		log.Logger.
			Info().
			Err(err).
			Msg("failed to listen and serve http server")
	}
}
