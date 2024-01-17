package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	PollCount      int64
	RandomValue    float64 = 123.0
	currentMetrics []entity.Metrics
)

func Run() error {
	// Initialize Agent Logger
	if err := logger.InitALogger(); err != nil {
		return err
	}

	// Parse agent configs
	var cfg Config
	Parse(&cfg)

	// Check server health
	if !serverIsHealthy(cfg.ServerAddress) {
		err := errors.New("please make sure that the server is up and running ")
		return err
	}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			pollMetrics()
		case <-reportTicker.C:
			reportMetrics(cfg.ServerAddress)
		}
	}
}

// floatPtr converts float to float pointer
func floatPtr(f float64) *float64 {
	return &f
}

// pollMetrics polls metrics from host that agent is running
func pollMetrics() {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	PollCount++
	currentMetrics = []entity.Metrics{
		{MType: entity.Counter, ID: "PollCount", Delta: &PollCount},
		{MType: entity.Gauge, ID: "RandomValue", Value: floatPtr(RandomValue)},
		{MType: entity.Gauge, ID: "Alloc", Value: floatPtr(float64(memStat.Alloc))},
		{MType: entity.Gauge, ID: "BuckHashSys", Value: floatPtr(float64(memStat.BuckHashSys))},
		{MType: entity.Gauge, ID: "Frees", Value: floatPtr(float64(memStat.Frees))},
		{MType: entity.Gauge, ID: "GCCPUFraction", Value: floatPtr(memStat.GCCPUFraction)},
		{MType: entity.Gauge, ID: "GCSys", Value: floatPtr(float64(memStat.GCSys))},
		{MType: entity.Gauge, ID: "HeapAlloc", Value: floatPtr(float64(memStat.HeapAlloc))},
		{MType: entity.Gauge, ID: "HeapIdle", Value: floatPtr(float64(memStat.HeapIdle))},
		{MType: entity.Gauge, ID: "HeapInuse", Value: floatPtr(float64(memStat.HeapInuse))},
		{MType: entity.Gauge, ID: "HeapObjects", Value: floatPtr(float64(memStat.HeapObjects))},
		{MType: entity.Gauge, ID: "HeapReleased", Value: floatPtr(float64(memStat.HeapReleased))},
		{MType: entity.Gauge, ID: "HeapSys", Value: floatPtr(float64(memStat.HeapSys))},
		{MType: entity.Gauge, ID: "LastGC", Value: floatPtr(float64(memStat.LastGC))},
		{MType: entity.Gauge, ID: "Lookups", Value: floatPtr(float64(memStat.Lookups))},
		{MType: entity.Gauge, ID: "MCacheInuse", Value: floatPtr(float64(memStat.MCacheInuse))},
		{MType: entity.Gauge, ID: "MCacheSys", Value: floatPtr(float64(memStat.MCacheSys))},
		{MType: entity.Gauge, ID: "MSpanInuse", Value: floatPtr(float64(memStat.MSpanInuse))},
		{MType: entity.Gauge, ID: "MSpanSys", Value: floatPtr(float64(memStat.MSpanSys))},
		{MType: entity.Gauge, ID: "Mallocs", Value: floatPtr(float64(memStat.Mallocs))},
		{MType: entity.Gauge, ID: "NextGC", Value: floatPtr(float64(memStat.NextGC))},
		{MType: entity.Gauge, ID: "NumForcedGC", Value: floatPtr(float64(memStat.NumForcedGC))},
		{MType: entity.Gauge, ID: "NumGC", Value: floatPtr(float64(memStat.NumGC))},
		{MType: entity.Gauge, ID: "OtherSys", Value: floatPtr(float64(memStat.OtherSys))},
		{MType: entity.Gauge, ID: "PauseTotalNs", Value: floatPtr(float64(memStat.PauseTotalNs))},
		{MType: entity.Gauge, ID: "StackInuse", Value: floatPtr(float64(memStat.StackInuse))},
		{MType: entity.Gauge, ID: "StackSys", Value: floatPtr(float64(memStat.StackSys))},
		{MType: entity.Gauge, ID: "Sys", Value: floatPtr(float64(memStat.Sys))},
		{MType: entity.Gauge, ID: "TotalAlloc", Value: floatPtr(float64(memStat.TotalAlloc))},
	}
}

// reportMetrics sends agent metrics to server
func reportMetrics(serverAddress string) {
	if len(currentMetrics) == 0 {
		logger.ALog.Info("No metrics to report.")
		return
	}
	for _, metric := range currentMetrics {
		url := fmt.Sprintf("http://%s/update", serverAddress)

		body, err := json.Marshal(metric)
		if err != nil {
			logger.ALog.Error("Error in reporting metrics:", err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			logger.ALog.Error("Error in reporting metrics:", err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			logger.ALog.Error("Error in closing response body", err)
			continue
		}

		logger.ALog.Info("Metrics reported successfully.")
	}
	currentMetrics = nil
}

// serverIsHealthy checks server health
func serverIsHealthy(serverAddress string) bool {
	url := fmt.Sprintf("http://%s/healthz", serverAddress)
	retry, delay := 3, 3*time.Second

	for i := 0; i < retry; i++ {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			logger.ALog.Error(err)
			log.Printf("Attempt #%d. Connection refused. Left %d attempts.\n", i+1, retry-i-1)
			time.Sleep(delay)
			continue
		}
		resp.Body.Close()
		return true
	}
	return false
}
