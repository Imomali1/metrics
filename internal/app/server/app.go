package server

import (
	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/services"
	"log"
	"net/http"
)

func Run() {
	memStorage := storage.NewStorage()
	repo := repository.New(memStorage)
	service := services.New(repo)
	handler := api.NewRouter(api.Options{
		ServiceManager: service,
	})
	log.Fatal(http.ListenAndServe(":8080", handler))
}
