package agent

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"crypto/rsa"

	"github.com/rs/zerolog/log"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

type agent struct {
	cfg       Config
	log       logger.Logger
	publicKey *rsa.PublicKey
}

type Metrics struct {
	mu        sync.RWMutex
	PollCount int64
	Arr       []entity.Metrics
}

func Run(cfg Config) {
	app := agent{
		cfg: cfg,
		log: logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName),
	}

	var err error
	if cfg.PublicKeyPath != "" {
		app.publicKey, err = cipher.UploadRSAPublicKey(cfg.PublicKeyPath)
		if err != nil {
			log.Err(err).Send()
			return
		}
	}

	if err = checkServer(cfg.ServerAddress); err != nil {
		log.Err(err).Send()
		return
	}

	log.Info().Msg("agent is up and running...")

	tasks := make(chan ReportTask)

	for i := 0; i < cfg.RateLimit; i++ {
		go app.worker(tasks)
	}

	metrics := new(Metrics)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	wg := &sync.WaitGroup{}

	app.pollRuntimeMetrics(metrics, wg, quit)
	app.pollGopsutilMetrics(metrics, wg, quit)
	app.reportMetricsV1(metrics, tasks, wg, quit)
	app.reportMetricsV2(metrics, tasks, wg, quit)
	app.reportMetricsV3(metrics, tasks, wg, quit)

	wg.Wait()
	close(tasks)

	log.Info().Msg("agent stopped successfully")
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
