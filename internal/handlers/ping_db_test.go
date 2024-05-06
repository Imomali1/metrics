package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMetricHandler_PingDB(t *testing.T) {
	r, mockServiceManager, handler := NewTestHandler()

	r.GET("/ping", handler.PingDB)

	tests := []struct {
		name         string
		mockSetup    func()
		expectedCode int
	}{
		{
			name: "successful ping",
			mockSetup: func() {
				mockServiceManager.On("Ping", mock.Anything).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "failed ping",
			mockSetup: func() {
				mockServiceManager.On("Ping", mock.Anything).Return(errors.New("db connection error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServiceManager.ExpectedCalls = nil
			tc.mockSetup()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)
			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
