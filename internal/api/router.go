package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
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
	gin.SetMode(gin.DebugMode)
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

	router.POST("/update/gauge/:name/:value", h.gaugeHandler.UpdateGauge)
	router.POST("/update/counter/:name/:value", h.counterHandler.UpdateCounter)

	return router
}
