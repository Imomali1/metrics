package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
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
}

func NewRouter(options Options) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.Use(middlewares.Logger())

	router.LoadHTMLGlob("static/templates/*.html")

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.ServiceManager.MetricService),
	}

	router.GET("/", h.MetricHandler.ListMetrics)
	router.POST("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.POST("/update/:type/:name/:value", h.MetricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", h.MetricHandler.GetMetricValueByName)

	return router
}
