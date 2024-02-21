package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	MetricHandler *handlers.MetricHandler
}

type Options struct {
	Logger         logger.Logger
	ServiceManager *services.Services
}

func NewRouter(options Options) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.ReqRespLogger(options.Logger))
	router.Use(middlewares.CompressResponse(), middlewares.DecompressRequest())

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.Logger, options.ServiceManager.MetricService),
	}

	router.LoadHTMLGlob("static/templates/*.html")

	router.GET("/", h.MetricHandler.ListMetrics)

	updateRoutes := router.Group("/update")
	{
		// v1 update handler using URI
		updateRoutes.POST("/:type/:name/:value", h.MetricHandler.UpdateMetricValue)

		// v2 update handler using JSON
		updateRoutes.POST("/", h.MetricHandler.UpdateMetricValueJSON)
	}

	getValueRoutes := router.Group("/value")
	{
		// v1 get value handler using URI
		getValueRoutes.GET("/:type/:name", h.MetricHandler.GetMetricValueByName)

		// v2 get value handler using JSON
		getValueRoutes.POST("/", h.MetricHandler.GetMetricValueByNameJSON)
	}

	return router
}
