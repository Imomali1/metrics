package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"log"
	"net/http"
)

func newHandler() http.Handler {
	memStorage := storage.New()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		ServiceManager: service,
	})
	return handler
}

func Run() {
	var cfg Config
	Parse(&cfg)

	handler := newHandler()
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler))
}
