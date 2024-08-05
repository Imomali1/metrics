package api

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"

	"crypto/rsa"

	"github.com/Imomali1/metrics/internal/handlers"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/usecase"
)

type Handlers struct {
	MetricHandler *handlers.MetricHandler
}

type Config struct {
	HashKey string
}

type Options struct {
	Logger           logger.Logger
	UseCase          usecase.UseCase
	Cfg              Config
	HTMLTemplatePath string
	PrivateKey       *rsa.PrivateKey
}

func NewRouter(options Options) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.ReqRespLogger(options.Logger))
	router.Use(middlewares.CompressResponse(), middlewares.DecompressRequest())
	router.Use(middlewares.ValidateHash(options.Logger, options.Cfg.HashKey))
	router.Use(middlewares.RSADecrypt(options.Logger, options.PrivateKey))

	h := Handlers{
		MetricHandler: handlers.NewMetricHandler(options.Logger, options.UseCase),
	}

	router.LoadHTMLGlob(options.HTMLTemplatePath)

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
