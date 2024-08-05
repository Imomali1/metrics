package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func ValidateHash(l logger.Logger, key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if key == "" {
			return
		}

		clientHash := ctx.GetHeader("HashSHA256")

		if clientHash != "" {
			data, err := utils.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				l.Logger.Info().Err(err).Msg("could not read body")
				return
			}

			controlHash := utils.GenerateHash(data, key)

			if clientHash != controlHash {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "request data not validated"})
				l.Logger.Info().Err(err).Msg("request data not validated")
				return
			}

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		ctx.Next()

		responseBody := w.body.Bytes()
		responseHash := utils.GenerateHash(responseBody, key)
		ctx.Header("HashSHA256", responseHash)
	}
}
