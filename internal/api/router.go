package api

import (
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/services"
	"github.com/gin-gonic/gin"
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

	router.LoadHTMLGlob("static/templates/*.html")

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.ServiceManager.MetricService),
	}

	router.GET("/", h.MetricHandler.ListMetrics)
	router.POST("/update/:type/:name/:value", h.MetricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", h.MetricHandler.GetMetricValueByName)

	return router
}
