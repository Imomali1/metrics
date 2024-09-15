package agent

import (
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/Imomali1/metrics/internal/entity"
)

type poller struct {
	interval time.Duration
}

func (a *agent) PollMetricsPeriodically(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(a.poller.interval)

	for {
		select {
		case <-ticker.C:
			go a.pollRuntimeMetrics(wg)
			go a.pollGopsutilMetrics(wg)
		case <-a.shutdownCh:
			a.log.Info().Msg("stopped collecting metrics")
			go a.reportMetricsV1(wg)
			go a.reportMetricsV2(wg)
			go a.reportMetricsV3(wg)
			return
		}
	}
}

func (a *agent) pollRuntimeMetrics(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	a.log.Info().Msg("started collecting runtime metrics")

	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	a.metrics.mu.Lock()

	a.metrics.PollCount++

	randomValue := rand.NormFloat64()

	a.metrics.Arr = []entity.Metrics{
		{MType: entity.Counter, ID: "PollCount", Delta: &a.metrics.PollCount},
		{MType: entity.Gauge, ID: "RandomValue", Value: utils.Ptr(randomValue)},
		{MType: entity.Gauge, ID: "Alloc", Value: utils.Ptr(float64(memStat.Alloc))},
		{MType: entity.Gauge, ID: "BuckHashSys", Value: utils.Ptr(float64(memStat.BuckHashSys))},
		{MType: entity.Gauge, ID: "Frees", Value: utils.Ptr(float64(memStat.Frees))},
		{MType: entity.Gauge, ID: "GCCPUFraction", Value: utils.Ptr(memStat.GCCPUFraction)},
		{MType: entity.Gauge, ID: "GCSys", Value: utils.Ptr(float64(memStat.GCSys))},
		{MType: entity.Gauge, ID: "HeapAlloc", Value: utils.Ptr(float64(memStat.HeapAlloc))},
		{MType: entity.Gauge, ID: "HeapIdle", Value: utils.Ptr(float64(memStat.HeapIdle))},
		{MType: entity.Gauge, ID: "HeapInuse", Value: utils.Ptr(float64(memStat.HeapInuse))},
		{MType: entity.Gauge, ID: "HeapObjects", Value: utils.Ptr(float64(memStat.HeapObjects))},
		{MType: entity.Gauge, ID: "HeapReleased", Value: utils.Ptr(float64(memStat.HeapReleased))},
		{MType: entity.Gauge, ID: "HeapSys", Value: utils.Ptr(float64(memStat.HeapSys))},
		{MType: entity.Gauge, ID: "LastGC", Value: utils.Ptr(float64(memStat.LastGC))},
		{MType: entity.Gauge, ID: "Lookups", Value: utils.Ptr(float64(memStat.Lookups))},
		{MType: entity.Gauge, ID: "MCacheInuse", Value: utils.Ptr(float64(memStat.MCacheInuse))},
		{MType: entity.Gauge, ID: "MCacheSys", Value: utils.Ptr(float64(memStat.MCacheSys))},
		{MType: entity.Gauge, ID: "MSpanInuse", Value: utils.Ptr(float64(memStat.MSpanInuse))},
		{MType: entity.Gauge, ID: "MSpanSys", Value: utils.Ptr(float64(memStat.MSpanSys))},
		{MType: entity.Gauge, ID: "Mallocs", Value: utils.Ptr(float64(memStat.Mallocs))},
		{MType: entity.Gauge, ID: "NextGC", Value: utils.Ptr(float64(memStat.NextGC))},
		{MType: entity.Gauge, ID: "NumForcedGC", Value: utils.Ptr(float64(memStat.NumForcedGC))},
		{MType: entity.Gauge, ID: "NumGC", Value: utils.Ptr(float64(memStat.NumGC))},
		{MType: entity.Gauge, ID: "OtherSys", Value: utils.Ptr(float64(memStat.OtherSys))},
		{MType: entity.Gauge, ID: "PauseTotalNs", Value: utils.Ptr(float64(memStat.PauseTotalNs))},
		{MType: entity.Gauge, ID: "StackInuse", Value: utils.Ptr(float64(memStat.StackInuse))},
		{MType: entity.Gauge, ID: "StackSys", Value: utils.Ptr(float64(memStat.StackSys))},
		{MType: entity.Gauge, ID: "Sys", Value: utils.Ptr(float64(memStat.Sys))},
		{MType: entity.Gauge, ID: "TotalAlloc", Value: utils.Ptr(float64(memStat.TotalAlloc))},
	}

	a.metrics.mu.Unlock()

	a.log.Info().Msg("finished collecting runtime metrics")
}

func (a *agent) pollGopsutilMetrics(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	a.log.Info().Msg("started collecting gopsutil metrics")

	vm, err := mem.VirtualMemory()
	if err != nil {
		a.log.Info().Err(err).Msg("cannot get memory metrics")
		return
	}

	total, free := float64(vm.Total), float64(vm.Free)

	a.metrics.mu.Lock()
	a.metrics.Arr = append(a.metrics.Arr, entity.Metrics{ID: "TotalMemory", MType: entity.Gauge, Value: &total})
	a.metrics.Arr = append(a.metrics.Arr, entity.Metrics{ID: "FreeMemory", MType: entity.Gauge, Value: &free})
	a.metrics.mu.Unlock()

	cpuUtil, err := cpu.Percent(0, false)
	if err != nil || len(cpuUtil) == 0 {
		a.log.Info().Err(err).Msg("cannot get cpu metrics")
		return
	}

	a.metrics.mu.Lock()
	a.metrics.Arr = append(a.metrics.Arr, entity.Metrics{ID: "CPUutilization1", MType: entity.Gauge, Value: &cpuUtil[0]})
	a.metrics.mu.Unlock()

	a.log.Info().Msg("finished collecting gopsutil metrics")
}
