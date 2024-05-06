package api

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/app/server/configs"
	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/services"
)

type Handlers struct {
	MetricHandler *handlers.MetricHandler
}

type Options struct {
	Logger         logger.Logger
	ServiceManager *services.Services
	Conf           configs.Config
}

func NewRouter(options Options) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.ReqRespLogger(options.Logger))
	router.Use(middlewares.CompressResponse(), middlewares.DecompressRequest())
	router.Use(middlewares.ValidateHash(options.Logger, options.Conf.HashKey))

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.Logger, options.ServiceManager),
	}

	router.LoadHTMLGlob("static/templates/*.html")

	router.GET("/", h.MetricHandler.ListMetrics)

	router.GET("/ping", h.MetricHandler.PingDB)

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	updateRoutes := router.Group("/update")
	{
		// v1 update handler using URI
		updateRoutes.POST("/:type/:name/:value", h.MetricHandler.UpdateMetricValue)

		// v2 update handler using JSON
		updateRoutes.POST("/", h.MetricHandler.UpdateMetricValueJSON)
	}

	updatesRoute := router.Group("/updates")
	{
		updatesRoute.POST("/", h.MetricHandler.Updates)
	}

	getValueRoutes := router.Group("/value")
	{
		// v1 get value handler using URI
		getValueRoutes.GET("/:type/:name", h.MetricHandler.GetMetricValueByName)

		// v2 get value handler using JSON
		getValueRoutes.POST("/", h.MetricHandler.GetMetricValueByNameJSON)
	}

	router.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	router.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	router.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	router.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	router.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))

	return router
}
