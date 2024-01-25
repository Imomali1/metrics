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

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.Logger, options.ServiceManager.MetricService),
	}

	router.LoadHTMLGlob("static/templates/*.html")

	router.GET("/", h.MetricHandler.ListMetrics)
	router.POST("/update/:type/:name/:value", h.MetricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", h.MetricHandler.GetMetricValueByName)
	router.POST("/update", h.MetricHandler.UpdateMetricValueJSON)
	router.GET("/value", h.MetricHandler.GetMetricValueByNameJSON)

	return router
}
