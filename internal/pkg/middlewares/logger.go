package middlewares

import (
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		uri := ctx.Request.RequestURI
		method := ctx.Request.Method

		ctx.Next()

		duration := time.Since(start)

		logger.Log.Infow("request",
			"uri", uri,
			"method", method,
			"duration", duration)

		status := ctx.Writer.Status()
		size := ctx.Writer.Size()

		logger.Log.Infow("response",
			"status", status,
			"size", size)
	}
}
