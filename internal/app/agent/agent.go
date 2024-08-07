package agent

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

type Metrics struct {
	mu        sync.RWMutex
	PollCount int64
	Arr       []entity.Metrics
}

func Run(cfg Config) {
	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	if err := checkServer(cfg.ServerAddress); err != nil {
		log.Err(err).Send()
		return
	}

	log.Info().Msg("agent is up and running...")

	tasks := make(chan ReportTask)

	for i := 0; i < cfg.RateLimit; i++ {
		go worker(log, tasks)
	}

	metrics := new(Metrics)

	var wg sync.WaitGroup
	wg.Add(5)
	go pollRuntimeMetrics(log, metrics, cfg.PollInterval, &wg)
	go pollGopsutilMetrics(log, metrics, cfg.PollInterval, &wg)
	go reportMetricsV1(log, cfg, metrics, tasks, &wg)
	go reportMetricsV2(log, cfg, metrics, tasks, &wg)
	go reportMetricsV3(log, cfg, metrics, tasks, &wg)
	wg.Wait()

	close(tasks)
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
