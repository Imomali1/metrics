package middlewares

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func CompressResponse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		acceptEncoding := ctx.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			ctx.Next()
			return
		}

		gz := gzip.NewWriter(ctx.Writer)
		defer gz.Close()

		ctx.Header("Content-Encoding", "gzip")
		ctx.Writer = &gzipWriter{
			ResponseWriter: ctx.Writer,
			writer:         gz,
		}
		ctx.Next()
	}
}
