package agent

import (
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	PollCount      int64
	RandomValue    float64 = 123.0
	currentMetrics []entity.Metric
)

func Run() {
	var cfg Config
	Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	log.Logger.Info().Msg("agent is up and running...")

	for {
		select {
		case <-pollTicker.C:
			log.Logger.Info().Msg("polling metrics...")
			pollMetrics()
		case <-reportTicker.C:
			log.Logger.Info().Msg("reporting metrics to server...")
			reportMetrics(log, cfg.ServerAddress)
		}
	}
}

func pollMetrics() {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	PollCount++
	currentMetrics = []entity.Metric{
		{Type: entity.Counter, Name: "PollCount", ValueCounter: PollCount},
		{Type: entity.Gauge, Name: "RandomValue", ValueGauge: RandomValue},
		{Type: entity.Gauge, Name: "Alloc", ValueGauge: float64(memStat.Alloc)},
		{Type: entity.Gauge, Name: "BuckHashSys", ValueGauge: float64(memStat.BuckHashSys)},
		{Type: entity.Gauge, Name: "Frees", ValueGauge: float64(memStat.Frees)},
		{Type: entity.Gauge, Name: "GCCPUFraction", ValueGauge: memStat.GCCPUFraction},
		{Type: entity.Gauge, Name: "GCSys", ValueGauge: float64(memStat.GCSys)},
		{Type: entity.Gauge, Name: "HeapAlloc", ValueGauge: float64(memStat.HeapAlloc)},
		{Type: entity.Gauge, Name: "HeapIdle", ValueGauge: float64(memStat.HeapIdle)},
		{Type: entity.Gauge, Name: "HeapInuse", ValueGauge: float64(memStat.HeapInuse)},
		{Type: entity.Gauge, Name: "HeapObjects", ValueGauge: float64(memStat.HeapObjects)},
		{Type: entity.Gauge, Name: "HeapReleased", ValueGauge: float64(memStat.HeapReleased)},
		{Type: entity.Gauge, Name: "HeapSys", ValueGauge: float64(memStat.HeapSys)},
		{Type: entity.Gauge, Name: "LastGC", ValueGauge: float64(memStat.LastGC)},
		{Type: entity.Gauge, Name: "Lookups", ValueGauge: float64(memStat.Lookups)},
		{Type: entity.Gauge, Name: "MCacheInuse", ValueGauge: float64(memStat.MCacheInuse)},
		{Type: entity.Gauge, Name: "MCacheSys", ValueGauge: float64(memStat.MCacheSys)},
		{Type: entity.Gauge, Name: "MSpanInuse", ValueGauge: float64(memStat.MSpanInuse)},
		{Type: entity.Gauge, Name: "MSpanSys", ValueGauge: float64(memStat.MSpanSys)},
		{Type: entity.Gauge, Name: "Mallocs", ValueGauge: float64(memStat.Mallocs)},
		{Type: entity.Gauge, Name: "NextGC", ValueGauge: float64(memStat.NextGC)},
		{Type: entity.Gauge, Name: "NumForcedGC", ValueGauge: float64(memStat.NumForcedGC)},
		{Type: entity.Gauge, Name: "NumGC", ValueGauge: float64(memStat.NumGC)},
		{Type: entity.Gauge, Name: "OtherSys", ValueGauge: float64(memStat.OtherSys)},
		{Type: entity.Gauge, Name: "PauseTotalNs", ValueGauge: float64(memStat.PauseTotalNs)},
		{Type: entity.Gauge, Name: "StackInuse", ValueGauge: float64(memStat.StackInuse)},
		{Type: entity.Gauge, Name: "StackSys", ValueGauge: float64(memStat.StackSys)},
		{Type: entity.Gauge, Name: "Sys", ValueGauge: float64(memStat.Sys)},
		{Type: entity.Gauge, Name: "TotalAlloc", ValueGauge: float64(memStat.TotalAlloc)},
	}
}

func reportMetrics(log logger.Logger, serverAddress string) {
	if len(currentMetrics) == 0 {
		log.Logger.Info().Msg("no metrics to report")
		return
	}
	for _, metric := range currentMetrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/", serverAddress, metric.Type, metric.Name)
		switch metric.Type {
		case entity.Counter:
			url = fmt.Sprintf("%s%d", url, metric.ValueCounter)
		case entity.Gauge:
			url = fmt.Sprintf("%s%f", url, metric.ValueGauge)
		default:
			log.Logger.Info().Msgf("invalid metric type: %s", metric.Type)
			continue
		}
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			log.Logger.Info().Err(err).Msg("error in reporting metrics")
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			log.Logger.Info().Err(err).Msg("error in closing response body")
			continue
		}

		log.Logger.Info().Msg("metrics reported successfully")
	}
	currentMetrics = nil
}
