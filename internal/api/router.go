package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	gaugeHandler   *handlers.GaugeHandler
	counterHandler *handlers.CounterHandler
}

type Options struct {
	ServiceManager *services.Services
}

func NewRouter(options Options) *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to metrics application!")
	})

	h := Handlers{
		gaugeHandler: handlers.NewGaugeHandler(handlers.GaugeHandlerOptions{
			ServiceManager: options.ServiceManager.GaugeService,
		}),
		counterHandler: handlers.NewCounterHandler(handlers.CounterHandlerOptions{
			ServiceManager: options.ServiceManager.CounterService,
		}),
	}

	metricsGroup := router.Group("/update")
	metricsGroup.Use(middlewares.ValidateUpdateURL())
	{
		// valid urls
		metricsGroup.POST("/gauge/:name/:value", h.gaugeHandler.UpdateGauge)
		metricsGroup.POST("/counter/:name/:value", h.counterHandler.UpdateCounter)

		// corner case invalid urls
		metricsGroup.POST("/gauge")
		metricsGroup.POST("/counter")
		metricsGroup.POST("/:type")
	}

	return router
}
