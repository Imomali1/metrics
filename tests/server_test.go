package tests

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/api"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/repository"
	"github.com/Imomali1/metrics/internal/usecase"
)

func setupRouter() *gin.Engine {
	store, _ := storage.New(context.Background(), "")
	repo := repository.New(store, nil)
	uc := usecase.New(repo)
	handler := api.NewRouter(api.Options{
		Logger:           logger.NewLogger(os.Stdout, "info", "test"),
		UseCase:          uc,
		Cfg:              api.Config{HashKey: "testKey"},
		HTMLTemplatePath: "../static/templates/*.html",
	})
	return handler
}

func TestServer(t *testing.T) {
	handler := setupRouter()

	tests := []struct {
		name               string
		method             string
		url                string
		requestContentType string
		requestBody        io.Reader
		wantedCode         int
		wantedBody         string
	}{
		/*================= PingDB =================*/
		{
			name:       "PingDB: ping memory storage",
			method:     http.MethodGet,
			url:        "/ping",
			wantedCode: http.StatusInternalServerError,
		},
		/*================= UpdateMetricValue =================*/
		{
			name:       "UpdateMetricValue: invalid metric type",
			method:     http.MethodPost,
			url:        "/update/invalid-type/name/123",
			wantedCode: http.StatusBadRequest,
		},
		{
			name:       "UpdateMetricValue: empty metric name",
			method:     http.MethodPost,
			url:        "/update/counter//123",
			wantedCode: http.StatusNotFound,
		},
		{
			name:       "UpdateMetricValue: empty metric name",
			method:     http.MethodPost,
			url:        "/update/gauge//123",
			wantedCode: http.StatusNotFound,
		},
		{
			name:       "UpdateMetricValue: invalid gauge value",
			method:     http.MethodPost,
			url:        "/update/gauge/gauge1/invalid-value",
			wantedCode: http.StatusBadRequest,
		},
		{
			name:       "UpdateMetricValue: invalid counter value",
			method:     http.MethodPost,
			url:        "/update/counter/counter1/invalid-value",
			wantedCode: http.StatusBadRequest,
		},
		{
			name:       "UpdateMetricValue: invalid counter value",
			method:     http.MethodPost,
			url:        "/update/counter/counter1/invalid-value",
			wantedCode: http.StatusBadRequest,
		},
		{
			name:       "UpdateMetricValue: valid gauge value",
			method:     http.MethodPost,
			url:        "/update/gauge/gauge1/123.4",
			wantedCode: http.StatusOK,
		},
		{
			name:       "UpdateMetricValue: valid counter value",
			method:     http.MethodPost,
			url:        "/update/counter/counter1/123",
			wantedCode: http.StatusOK,
		},
		/*================= GetMetricValueByName =================*/
		{
			name:       "GetMetricValueByName: invalid metric type",
			method:     http.MethodGet,
			url:        "/value/invalid-type/name",
			wantedCode: http.StatusBadRequest,
		},
		{
			name:       "GetMetricValueByName: non-existing counter",
			url:        "/value/counter/name",
			wantedCode: http.StatusNotFound,
		},
		{
			name:       "GetMetricValueByName: non-existing gauge",
			url:        "/value/gauge/name",
			wantedCode: http.StatusNotFound,
		},
		{
			name:       "GetMetricValueByName: existing counter",
			url:        "/value/counter/counter1",
			wantedCode: http.StatusOK,
			wantedBody: "123",
		},
		{
			name:       "GetMetricValueByName: existing gauge",
			url:        "/value/gauge/gauge1",
			wantedCode: http.StatusOK,
			wantedBody: "123.4",
		},
		/*================= UpdateMetricValueJSON =================*/
		{
			name:               "UpdateMetricValueJSON: invalid request content type",
			method:             http.MethodPost,
			url:                "/update/",
			requestContentType: "invalid-content-type",
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "UpdateMetricValueJSON: invalid body",
			method:             http.MethodPost,
			url:                "/update/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"invalid-key": "invalid-value"}`),
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "UpdateMetricValueJSON: invalid metric type",
			method:             http.MethodPost,
			url:                "/update/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"type": "invalid-type"}`),
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "UpdateMetricValueJSON: valid counter update",
			method:             http.MethodPost,
			url:                "/update/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id": "counter2","type":"counter","delta": 321}`),
			wantedCode:         http.StatusOK,
		},
		{
			name:               "UpdateMetricValueJSON: valid gauge update",
			method:             http.MethodPost,
			url:                "/update/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id": "gauge2","type":"gauge","value": 321.5}`),
			wantedCode:         http.StatusOK,
		},
		/*================= GetMetricValueByNameJSON =================*/
		{
			name:               "GetMetricValueByNameJSON: invalid request content type",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "invalid-content-type",
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "GetMetricValueByNameJSON: invalid body",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"invalid-key": "invalid-value"}`),
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "GetMetricValueByNameJSON: invalid metric type",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"type": "invalid-type"}`),
			wantedCode:         http.StatusBadRequest,
		},
		{
			name:               "GetMetricValueByNameJSON: non-existing counter",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id":"non-existing","type":"counter"}`),
			wantedCode:         http.StatusNotFound,
		},
		{
			name:               "GetMetricValueByNameJSON: non-existing gauge",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id": "non-existing","type":"gauge"}`),
			wantedCode:         http.StatusNotFound,
		},
		{
			name:               "GetMetricValueByNameJSON: get valid counter",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id": "counter2","type":"counter"}`),
			wantedCode:         http.StatusOK,
			wantedBody:         `{"id": "counter2","type":"counter","delta": 321}`,
		},
		{
			name:               "UpdateMetricValueJSON: get valid gauge",
			method:             http.MethodPost,
			url:                "/value/",
			requestContentType: "application/json",
			requestBody:        strings.NewReader(`{"id": "gauge2","type":"gauge"}`),
			wantedCode:         http.StatusOK,
			wantedBody:         `{"id": "gauge2","type":"gauge","value": 321.5}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, test.url, test.requestBody)
			require.NoError(t, err)

			request.Header.Set("Content-Type", test.requestContentType)

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			require.Equal(t, response.Code, test.wantedCode)

			if test.wantedBody == "" {
				return
			}

			if test.requestContentType == "application/json" {
				require.JSONEq(t, test.wantedBody, response.Body.String())
				return
			}

			require.Equal(t, test.wantedBody, response.Body.String())
		})
	}
}
