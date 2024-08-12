package agent

import (
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/Imomali1/metrics/internal/entity"
)

func (a agent) pollRuntimeMetrics(metrics *Metrics, wg *sync.WaitGroup, quit <-chan os.Signal) {
	ticker := time.NewTicker(a.cfg.PollInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()

				a.log.Info().Msg("started collecting runtime metrics")
				var memStat runtime.MemStats
				runtime.ReadMemStats(&memStat)
				metrics.mu.Lock()
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
				metrics.mu.Unlock()
				a.log.Info().Msg("finished collecting runtime metrics")
			}()
		case <-quit:
			a.log.Info().Msg("stopping collecting runtime metrics")
			return
		}
	}
}

func (a agent) pollGopsutilMetrics(metrics *Metrics, wg *sync.WaitGroup, quit <-chan os.Signal) {
	ticker := time.NewTicker(a.cfg.PollInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()

				a.log.Info().Msg("started collecting gopsutil metrics")
				vm, err := mem.VirtualMemory()
				if err != nil {
					a.log.Info().Err(err).Msg("cannot get memory metrics")
					return
				}

				total, free := float64(vm.Total), float64(vm.Free)

				metrics.mu.Lock()
				metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "TotalMemory", MType: entity.Gauge, Value: &total})
				metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "FreeMemory", MType: entity.Gauge, Value: &free})
				metrics.mu.Unlock()

				cpuUtil, err := cpu.Percent(0, false)
				if err != nil || len(cpuUtil) == 0 {
					a.log.Info().Err(err).Msg("cannot get cpu metrics")
					return
				}

				metrics.mu.Lock()
				metrics.Arr = append(metrics.Arr, entity.Metrics{ID: "CPUutilization1", MType: entity.Gauge, Value: &cpuUtil[0]})
				metrics.mu.Unlock()

				a.log.Info().Msg("finished collecting gopsutil metrics")
			}()
		case <-quit:
			a.log.Info().Msg("stopping collecting gopsutil metrics")
			return
		}
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
