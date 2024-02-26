package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	PollCount int64
	Arr       []entity.Metrics
	mu        sync.RWMutex
}

func Run() {
	var cfg Config
	Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	metrics := new(Metrics)

	log.Logger.Info().Msg("agent is up and running...")

	for {
		select {
		case <-pollTicker.C:
			log.Logger.Info().Msg("polling metrics...")
			pollMetrics(metrics)
		case <-reportTicker.C:
			log.Logger.Info().Msg("reporting metrics to server/v1...")
			reportMetricsV1(log, cfg.ServerAddress, metrics)
			log.Logger.Info().Msg("reporting metrics to server/v2...")
			reportMetricsV2(log, cfg.ServerAddress, metrics)
		}
	}
}

func pollMetrics(metrics *Metrics) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	metrics.PollCount++
	RandomValue := rand.NormFloat64()
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.Arr = []entity.Metrics{
		{MType: entity.Counter, ID: "PollCount", Delta: &metrics.PollCount},
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

func floatPtr(f float64) *float64 {
	return &f
}

func reportMetricsV1(log logger.Logger, serverAddress string, metrics *Metrics) {
	if len(metrics.Arr) == 0 {
		log.Logger.Info().Msg("no metrics to report")
		return
	}
	client := resty.New().SetHeader("Content-Type", "text/plain")

	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	for _, metric := range metrics.Arr {
		url := fmt.Sprintf("http://%s/update/%s/%s/", serverAddress, metric.MType, metric.ID)
		switch metric.MType {
		case entity.Counter:
			url = fmt.Sprintf("%s%d", url, *metric.Delta)
		case entity.Gauge:
			url = fmt.Sprintf("%s%f", url, *metric.Value)
		default:
			log.Logger.Info().Msgf("invalid metric type: %s", metric.MType)
			continue
		}
		_, err := client.R().Post(url)
		if err != nil {
			log.Logger.Info().Err(err).Msg("error in making request")
			continue
		}

		log.Logger.Info().Msg("metrics reported successfully")
	}
}

func reportMetricsV2(log logger.Logger, serverAddress string, metrics *Metrics) {
	if len(metrics.Arr) == 0 {
		log.Logger.Info().Msg("no metrics to report")
		return
	}
	client := resty.New().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")
	url := fmt.Sprintf("http://%s/update/", serverAddress)
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()
	for _, metric := range metrics.Arr {
		body, err := easyjson.Marshal(metric)
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot unmarshal metric object")
			continue
		}
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		_, err = gzipWriter.Write(body)
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot compress body")
			continue
		}
		err = gzipWriter.Close()
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot close gzip writer")
			continue
		}
		_, err = client.R().
			SetBody(buf.Bytes()).
			Post(url)
		if err != nil {
			log.Logger.Info().Err(err).Msg("error in making request")
			continue
		}

		log.Logger.Info().Msg("metrics reported successfully")
	}
}
