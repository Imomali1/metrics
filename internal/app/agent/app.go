package agent

import (
	"flag"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	Gauge   entity.MetricType = "gauge"
	Counter entity.MetricType = "counter"
)

var (
	PollCount      int64
	RandomValue    float64 = 123.0
	currentMetrics []entity.Metric
)

func Run() {
	serverAddress := flag.String("a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", 2, "частота отправки метрик на сервер")
	reportInterval := flag.Int("r", 10, "частота опроса метрик из пакета runtime")
	flag.Parse()

	pollTicker := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(*reportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			pollMetrics()
		case <-reportTicker.C:
			reportMetrics(*serverAddress)
		}
	}
}

func pollMetrics() {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	PollCount++
	currentMetrics = []entity.Metric{
		{Type: Counter, Name: "PollCount", Value: PollCount},
		{Type: Gauge, Name: "RandomValue", Value: RandomValue},
		{Type: Gauge, Name: "Alloc", Value: float64(memStat.Alloc)},
		{Type: Gauge, Name: "BuckHashSys", Value: float64(memStat.BuckHashSys)},
		{Type: Gauge, Name: "Frees", Value: float64(memStat.Frees)},
		{Type: Gauge, Name: "GCCPUFraction", Value: memStat.GCCPUFraction},
		{Type: Gauge, Name: "GCSys", Value: float64(memStat.GCSys)},
		{Type: Gauge, Name: "HeapAlloc", Value: float64(memStat.HeapAlloc)},
		{Type: Gauge, Name: "HeapIdle", Value: float64(memStat.HeapIdle)},
		{Type: Gauge, Name: "HeapInuse", Value: float64(memStat.HeapInuse)},
		{Type: Gauge, Name: "HeapObjects", Value: float64(memStat.HeapObjects)},
		{Type: Gauge, Name: "HeapReleased", Value: float64(memStat.HeapReleased)},
		{Type: Gauge, Name: "HeapSys", Value: float64(memStat.HeapSys)},
		{Type: Gauge, Name: "LastGC", Value: float64(memStat.LastGC)},
		{Type: Gauge, Name: "Lookups", Value: float64(memStat.Lookups)},
		{Type: Gauge, Name: "MCacheInuse", Value: float64(memStat.MCacheInuse)},
		{Type: Gauge, Name: "MCacheSys", Value: float64(memStat.MCacheSys)},
		{Type: Gauge, Name: "MSpanInuse", Value: float64(memStat.MSpanInuse)},
		{Type: Gauge, Name: "MSpanSys", Value: float64(memStat.MSpanSys)},
		{Type: Gauge, Name: "Mallocs", Value: float64(memStat.Mallocs)},
		{Type: Gauge, Name: "NextGC", Value: float64(memStat.NextGC)},
		{Type: Gauge, Name: "NumForcedGC", Value: float64(memStat.NumForcedGC)},
		{Type: Gauge, Name: "NumGC", Value: float64(memStat.NumGC)},
		{Type: Gauge, Name: "OtherSys", Value: float64(memStat.OtherSys)},
		{Type: Gauge, Name: "PauseTotalNs", Value: float64(memStat.PauseTotalNs)},
		{Type: Gauge, Name: "StackInuse", Value: float64(memStat.StackInuse)},
		{Type: Gauge, Name: "StackSys", Value: float64(memStat.StackSys)},
		{Type: Gauge, Name: "Sys", Value: float64(memStat.Sys)},
		{Type: Gauge, Name: "TotalAlloc", Value: float64(memStat.TotalAlloc)},
	}
}

func reportMetrics(serverAddress string) {
	if len(currentMetrics) == 0 {
		log.Println("No metrics to report.")
		return
	}
	for _, metric := range currentMetrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v",
			serverAddress, metric.Type, metric.Name, metric.Value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Error in reporting metrics:", err)
			continue
		}
		defer resp.Body.Close()

		fmt.Println("Metrics reported successfully.")
	}
	currentMetrics = nil
}
