package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

type reporter struct {
	interval  time.Duration
	publicKey *rsa.PublicKey
}

func (a *agent) ReportMetricsPeriodically(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(a.reporter.interval)

	for {
		select {
		case <-ticker.C:
			go a.reportMetricsV1(wg)
			go a.reportMetricsV2(wg)
			go a.reportMetricsV3(wg)
		case <-a.shutdownCh:
			log.Info().Msg("stopped reporting metrics to server periodically")
			return
		}
	}
}

func (a *agent) reportMetricsV1(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	a.log.Info().Msg("started reporting metrics to server/v1...")
	if len(a.metrics.Arr) == 0 {
		a.log.Info().Msg("no metrics to report")
		return
	}

	client := resty.New().SetHeader("Content-Type", "text/plain")

	a.metrics.mu.RLock()
	arr := a.metrics.Arr
	a.metrics.mu.RUnlock()

	for _, metric := range arr {
		url := fmt.Sprintf("http://%s/update/%s/%s/", a.cfg.ServerAddress, metric.MType, metric.ID)
		switch metric.MType {
		case entity.Counter:
			url = fmt.Sprintf("%s%d", url, *metric.Delta)
		case entity.Gauge:
			url = fmt.Sprintf("%s%f", url, *metric.Value)
		default:
			a.log.Info().Msgf("invalid metric type: %s", metric.MType)
			continue
		}

		a.jobsChan <- Job{
			Request: client.R(),
			URL:     url,
		}
	}

	a.log.Info().Msg("finished reporting metrics to server/v1...")
}

func (a *agent) reportMetricsV2(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	a.log.Info().Msg("started reporting metrics to server/v2...")
	if len(a.metrics.Arr) == 0 {
		a.log.Info().Msg("no metrics to report")
		return
	}

	client := resty.New().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	url := fmt.Sprintf("http://%s/update/", a.cfg.ServerAddress)

	a.metrics.mu.RLock()
	arr := a.metrics.Arr
	a.metrics.mu.RUnlock()

	for _, metric := range arr {
		body, err := easyjson.Marshal(metric)
		if err != nil {
			a.log.Info().Err(err).Msg("cannot unmarshal metric object")
			continue
		}

		if a.reporter.publicKey != nil {
			body, err = cipher.EncryptRSA(a.reporter.publicKey, body)
			if err != nil {
				a.log.Info().Err(err).Msg("cannot encrypt message")
				continue
			}
		}

		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)

		_, err = gzipWriter.Write(body)
		if err != nil {
			a.log.Info().Err(err).Msg("cannot compress body")
			continue
		}

		err = gzipWriter.Close()
		if err != nil {
			a.log.Info().Err(err).Msg("cannot close gzip writer")
			continue
		}

		if a.cfg.HashKey != "" {
			hash := utils.GenerateHash(buf.Bytes(), a.cfg.HashKey)
			client.SetHeader("HashSHA256", hash)
		}

		a.jobsChan <- Job{
			Request: client.R().SetBody(buf.Bytes()),
			URL:     url,
		}
	}

	a.log.Info().Msg("finished reporting metrics to server/v2...")
}

func (a *agent) reportMetricsV3(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	a.log.Info().Msg("started reporting metrics to server/v3...")
	if len(a.metrics.Arr) == 0 {
		a.log.Info().Msg("no metrics to report")
		return
	}

	client := resty.New().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	url := fmt.Sprintf("http://%s/updates/", a.cfg.ServerAddress)

	a.metrics.mu.RLock()
	arr := a.metrics.Arr
	a.metrics.mu.RUnlock()

	list := entity.MetricsList(arr)
	body, err := easyjson.Marshal(&list)
	if err != nil {
		a.log.Info().Err(err).Msg("cannot unmarshal metric object")
		return
	}

	if a.reporter.publicKey != nil {
		body, err = cipher.EncryptRSA(a.reporter.publicKey, body)
		if err != nil {
			a.log.Info().Err(err).Msg("cannot encrypt message")
			return
		}
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(body)
	if err != nil {
		a.log.Info().Err(err).Msg("cannot compress body")
		return
	}

	err = gzipWriter.Close()
	if err != nil {
		a.log.Info().Err(err).Msg("cannot close gzip writer")
		return
	}

	if a.cfg.HashKey != "" {
		hash := utils.GenerateHash(buf.Bytes(), a.cfg.HashKey)
		client.SetHeader("HashSHA256", hash)
	}

	a.jobsChan <- Job{
		Request: client.R().SetBody(buf.Bytes()),
		URL:     url,
	}

	a.log.Info().Msg("finished reporting metrics to server/v3...")
}
