package middlewares

import (
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func ReqRespLogger(l logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		uri := ctx.Request.RequestURI
		method := ctx.Request.Method

		ctx.Next()

		l.Logger.
			Info().
			Str("uri", uri).
			Str("method", method).
			Str("duration", time.Since(start).String()).
			Msg("request details")

		status := ctx.Writer.Status()
		size := ctx.Writer.Size()

		l.Logger.
			Info().
			Int("status", status).
			Int("size", size).
			Msg("response details")
	}
}
