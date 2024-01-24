package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	MetricHandler *handlers.MetricHandler
}

type Options struct {
	ServiceManager *services.Services
	Logger         logger.Logger
}

func NewRouter(options Options) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestResponseLogger(options.Logger))

	router.LoadHTMLGlob("static/templates/*.html")

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.ServiceManager.MetricService),
	}

	router.GET("/", h.MetricHandler.ListMetrics)
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.POST("/update", h.MetricHandler.UpdateMetricValue)
	router.GET("/value", h.MetricHandler.GetMetricValueByName)

	return router
}
