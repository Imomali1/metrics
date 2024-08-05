package middlewares

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func RSADecrypt(l logger.Logger, privateKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if privateKey == nil {
			return
		}

		data, err := utils.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			l.Logger.Info().Err(err).Msg("could not read body")
			return
		}

		decryptedData, err := cipher.DecryptRSA(privateKey, data)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
			l.Logger.Info().Err(err).Msg("failed to decrypt data")
			return
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedData))

		ctx.Next()
	}
}
