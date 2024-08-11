package tests

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/middlewares"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func TestRSADecrypt(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	gin.SetMode(gin.TestMode)

	l := logger.NewLogger(os.Stdout, "info", "test")

	router := gin.New()

	router.Use(middlewares.RSADecrypt(l, privateKey))

	router.POST("/test", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	rawBody := []byte(`{"message": "test"}`)
	body, _ := cipher.EncryptRSA(&privateKey.PublicKey, rawBody)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	newReqBody, _ := io.ReadAll(req.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, rawBody, newReqBody)
}

func TestReqRespLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	l := logger.NewLogger(os.Stdout, "info", "test")

	router := gin.New()

	router.Use(middlewares.ReqRespLogger(l))

	router.POST("/test", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	body := []byte(`{"message": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateHash(t *testing.T) {
	gin.SetMode(gin.TestMode)

	l := logger.NewLogger(os.Stdout, "info", "test")

	tts := []struct {
		name string
		key  string
	}{
		{
			name: "empty key",
			key:  "",
		},
		{
			name: "normal key",
			key:  "test key",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.Use(middlewares.ValidateHash(l, tt.key))

			router.POST("/test", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			body := []byte(`{"message": "test"}`)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("HashSHA256", utils.GenerateHash(body, tt.key))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			if tt.key == "" {
				require.Empty(t, w.Header().Get("HashSHA256"))
				return
			}

			require.NotEmpty(t, w.Header().Get("HashSHA256"))
		})
	}
}
