package api

import (
	"context"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func SetupRouterOptions() Options {
	log := logger.NewLogger(os.Stdout, "info", "test")
	store, _ := storage.New(context.Background(), "")
	repo := repository.New(store, nil)
	uc := usecase.New(repo)

	return Options{
		Logger:  log,
		UseCase: uc,
		Cfg: Config{
			HashKey: "testKey",
		},
		HTMLTemplatePath: "../../static/templates/*.html",
	}
}

func SetupRouter() *gin.Engine {
	options := SetupRouterOptions()
	router := NewRouter(options)
	return router
}

func TestNewRouter(t *testing.T) {
	router := SetupRouter()
	assert.NotNil(t, router)
}
