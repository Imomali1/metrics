package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func newHandler(log logger.Logger) *gin.Engine {
	memStorage := storage.NewStorage()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:         log,
		ServiceManager: service,
	})
	return handler
}

func Run() {
	var cfg Config
	Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	handler := newHandler(log)

	log.Logger.Info().Msgf("server is up and listening on address: %s", cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, handler)
	if err != nil {
		log.Logger.
			Info().
			Err(err).
			Msg("failed to listen and serve http server")
	}
}
