package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func NTestServer(t *testing.T) {
	ts := httptest.NewServer(newHandler())
	defer ts.Close()

	var updateMetricTestCases = []struct {
		name   string
		url    string
		status int
	}{
		{"counter_without_id", "/update/counter", http.StatusNotFound},
		{"gauge_without_id", "/update/gauge", http.StatusNotFound},
		{"unknown_without_id", "/update/unknown", http.StatusNotFound},
		{"counter_invalid_value", "/update/counter/testCounter/invalid", http.StatusBadRequest},
		{"gauge_invalid_value", "/update/gauge/testGauge/invalid", http.StatusBadRequest},
		{"unknown_string_value", "/update/unknown/testUnknown/string", http.StatusBadRequest},
		{"unknown_int_value", "/update/unknown/testUnknown/12", http.StatusBadRequest},
		{"unknown_float_value", "/update/unknown/testUnknown/12.50", http.StatusBadRequest},
		{"counter_valid_update", "/update/counter/testCounter/12", http.StatusOK},
		{"gauge_valid_update", "/update/gauge/testGauge/12.50", http.StatusOK},
	}
	for _, v := range updateMetricTestCases {
		t.Run(v.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, http.MethodPost, v.url)
			assert.Equal(t, v.status, resp.StatusCode)
			err := resp.Body.Close()
			assert.NoError(t, err)
		})
	}

	var getMetricTestCases = []struct {
		name   string
		url    string
		want   string
		status int
	}{
		{"counter_not_found", "/value/counter/not_found", "", http.StatusNotFound},
		{"gauge_not_found", "/value/gauge/not_found", "", http.StatusNotFound},
		{"unknown_metric", "/value/unknown/any", "", http.StatusNotFound},
		{"counter_valid_value", "/value/counter/testCounter", "12", http.StatusOK},
		{"gauge_valid_value", "/value/gauge/testGauge", "12.5", http.StatusOK},
	}
	for _, v := range getMetricTestCases {
		t.Run(v.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, "GET", v.url)
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.want, get)
			err := resp.Body.Close()
			assert.NoError(t, err)
		})
	}
}
