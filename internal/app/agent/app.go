package agent

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/Imomali1/metrics/internal/app/agent/configs"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.RWMutex
	PollCount int64
	Arr       []entity.Metrics
}

func Run() {
	var cfg configs.Config
	configs.Parse(&cfg)

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	if err := checkServer(cfg.ServerAddress); err != nil {
		log.Logger.Err(err).Send()
		return
	}

	log.Logger.Info().Msg("agent is up and running...")

	metrics := new(Metrics)

	reportWorker := make(chan []entity.Metrics, cfg.RateLimit)
	go reportMetricsWorker(log, cfg.ServerAddress, cfg.HashKey, reportWorker)

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			log.Logger.Info().Msg("polling runtime metrics...")
			go pollRuntimeMetrics(metrics)
			log.Logger.Info().Msg("polling gopsutil metrics...")
			go pollGopsutilMetrics(log, metrics)
		case <-reportTicker.C:
			log.Logger.Info().Msg("queueing metrics for reporting...")
			reportWorker <- metrics.Arr
		}
	}
}

func pollRuntimeMetrics(metrics *Metrics) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	metrics.PollCount++
	randomValue := rand.NormFloat64()
	metrics.Arr = []entity.Metrics{
		{MType: entity.Counter, ID: "PollCount", Delta: &metrics.PollCount},
		{MType: entity.Gauge, ID: "RandomValue", Value: floatPtr(randomValue)},
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

func pollGopsutilMetrics(log logger.Logger, metrics *Metrics) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Logger.Info().Err(err).Msg("cannot get memory metrics")
		return
	}

	total, free := float64(vm.Total), float64(vm.Free)

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "TotalMemory", MType: entity.Gauge, Value: &total})
	metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "FreeMemory", MType: entity.Gauge, Value: &free})

	cpuUtil, err := cpu.Percent(0, false)
	if err != nil || len(cpuUtil) == 0 {
		log.Logger.Info().Err(err).Msg("cannot get cpu metrics")
		return
	}

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "CPUutilization1", MType: entity.Gauge, Value: &cpuUtil[0]})
}

func reportMetricsWorker(log logger.Logger, serverAddress string, hashKey string, reportWorker <-chan []entity.Metrics) {
	client := resty.New().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	url := fmt.Sprintf("http://%s/updates/", serverAddress)

	for metricsArr := range reportWorker {
		if len(metricsArr) == 0 {
			log.Logger.Info().Msg("no metrics to report")
			continue
		}

		batch := entity.MetricsList(metricsArr)
		body, err := easyjson.Marshal(&batch)
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

		if hashKey != "" {
			hash := utils.GenerateHash(buf.Bytes(), hashKey)
			client.SetHeader("HashSHA256", hash)
		}

		err = utils.DoWithRetries(func() error {
			_, err = client.R().
				SetBody(buf.Bytes()).
				Post(url)
			return err
		})

		if err != nil {
			log.Logger.Info().Err(err).Msg("error in making request")
			continue
		}

		log.Logger.Info().Msg("metrics reported successfully")
	}
}

func checkServer(address string) error {
	client := resty.New()
	url := fmt.Sprintf("http://%s/healthz", address)

	var err error
	err = utils.DoWithRetries(func() error {
		_, err = client.R().Get(url)
		return err
	})

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return opErr.Err
	}
	return err
}

func floatPtr(f float64) *float64 {
	return &f
}
