package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	MetricHandler *handlers.MetricHandler
}

type Options struct {
	ServiceManager *services.Services
}

func NewRouter(options Options) *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(handlers.MetricHandlerOptions{
			ServiceManager: options.ServiceManager,
		}),
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to metrics application!")
	})

	router.POST("/update/:type/:name/:value", h.MetricHandler.UpdateMetric)

	return router
}
