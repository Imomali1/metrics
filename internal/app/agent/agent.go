package agent

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"github.com/go-resty/resty/v2"
)

type agent struct {
	cfg        Config
	log        logger.Logger
	metrics    *Metrics
	poller     poller
	reporter   reporter
	jobsChan   chan Job
	shutdownCh chan struct{}
}

type Metrics struct {
	mu        sync.RWMutex
	PollCount int64
	Arr       []entity.Metrics
}

func Run(cfg Config, log logger.Logger) error {
	publicKey, err := cipher.UploadRSAPublicKey(cfg.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to upload public key: %w", err)
	}

	app := agent{
		cfg:     cfg,
		log:     log,
		metrics: &Metrics{},
		poller: poller{
			interval: time.Duration(cfg.PollInterval) * time.Second,
		},
		reporter: reporter{
			interval:  time.Duration(cfg.ReportInterval) * time.Second,
			publicKey: publicKey,
		},
		jobsChan:   make(chan Job),
		shutdownCh: make(chan struct{}),
	}

	if err = checkServer(cfg.ServerAddress); err != nil {
		return fmt.Errorf("failed to check server: %w", err)
	}

	log.Info().Msg("agent is up and running...")

	for i := 0; i < cfg.RateLimit; i++ {
		go app.worker()
	}

	var wg sync.WaitGroup

	go app.PollMetricsPeriodically(&wg)
	go app.ReportMetricsPeriodically(&wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM|syscall.SIGINT|syscall.SIGQUIT)

	<-quit
	close(app.shutdownCh)
	wg.Wait()
	close(app.jobsChan)

	return nil
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
