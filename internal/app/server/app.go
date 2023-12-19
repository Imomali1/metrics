package server

import (
	"fmt"
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"log"
	"net/http"
)

func Run() {
	var cfg Config
	Parse(&cfg)
	fmt.Println(cfg)

	memStorage := storage.NewStorage()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		ServiceManager: service,
	})
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler))
}
