package middlewares

import (
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

func RequestResponseLogger(l logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		uri := ctx.Request.RequestURI
		method := ctx.Request.Method

		ctx.Next()

		l.Logger.
			Info().
			Dict("request", zerolog.Dict().
				Str("uri", uri).
				Str("method", method).
				Str("duration", time.Since(start).String()),
			).Send()

		status := ctx.Writer.Status()
		size := ctx.Writer.Size()

		l.Logger.
			Info().
			Dict("response", zerolog.Dict().
				Int("status", status).
				Int("size", size),
			).Send()
	}
}
