package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPollMetrics(t *testing.T) {
	Run()
	time.Sleep(3 * time.Second)
	if len(currentMetrics) == 0 {
		t.Errorf("Expected updated metrics, but got empty slice!")
	}
	t.Log("Successfully polled!")
	t.Fatal("Test should be stopped due to a timeout.")
}

func TestReportMetrics(t *testing.T) {
	tempServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	currentMetrics = []Metric{
		{Type: Gauge, Name: "Test", Value: 123.0},
	}

	serverAddress = tempServer.URL

	reportMetrics()

	if len(currentMetrics) != 0 {
		t.Errorf("Expected metrics to be reset after sending to the server.")
	}
}
