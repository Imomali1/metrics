package middlewares

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type gzipReader struct {
	io.ReadCloser
	reader *gzip.Reader
}

func (g *gzipReader) Read(data []byte) (int, error) {
	return g.reader.Read(data)
}

func (g *gzipReader) Close() error {
	return g.reader.Close()
}

func DecompressRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !strings.Contains(ctx.GetHeader("Content-Encoding"), "gzip") {
			ctx.Next()
			return
		}
		gz, err := gzip.NewReader(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer gz.Close()

		ctx.Request.Body = &gzipReader{
			ReadCloser: ctx.Request.Body,
			reader:     gz,
		}
		ctx.Next()
	}
}
