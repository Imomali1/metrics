package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/Imomali1/metrics/internal/app/agent/configs"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"sync"
	"time"
)

type ReportTask struct {
	Request *resty.Request
	URL     string
}

func (t *ReportTask) Process() error {
	err := utils.DoWithRetries(func() error {
		_, err := t.Request.Post(t.URL)
		return err
	})
	return err
}

func worker(log logger.Logger, tasks <-chan ReportTask) {
	for task := range tasks {
		if err := task.Process(); err != nil {
			log.Logger.Info().Err(err).Msg("error in reporting metrics to server")
		} else {
			log.Logger.Info().Msg("metrics reported successfully")
		}
	}
}

func reportMetricsV1(log logger.Logger, cfg configs.Config, metrics *Metrics, requests chan<- ReportTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
		log.Logger.Info().Msg("started reporting metrics to server/v1...")
		if len(metrics.Arr) == 0 {
			log.Logger.Info().Msg("no metrics to report")
			return
		}

		client := resty.New().SetHeader("Content-Type", "text/plain")

		metrics.mu.RLock()
		arr := metrics.Arr
		metrics.mu.RUnlock()

		for _, metric := range arr {
			url := fmt.Sprintf("http://%s/update/%s/%s/", cfg.ServerAddress, metric.MType, metric.ID)
			switch metric.MType {
			case entity.Counter:
				url = fmt.Sprintf("%s%d", url, *metric.Delta)
			case entity.Gauge:
				url = fmt.Sprintf("%s%f", url, *metric.Value)
			default:
				log.Logger.Info().Msgf("invalid metric type: %s", metric.MType)
				continue
			}

			requests <- ReportTask{
				Request: client.R(),
				URL:     url,
			}
		}
		log.Logger.Info().Msg("finished reporting metrics to server/v1...")
	}
}

func reportMetricsV2(log logger.Logger, cfg configs.Config, metrics *Metrics, requests chan<- ReportTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
		log.Logger.Info().Msg("started reporting metrics to server/v2...")
		if len(metrics.Arr) == 0 {
			log.Logger.Info().Msg("no metrics to report")
			return
		}

		client := resty.New().
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Content-Type", "application/json")

		url := fmt.Sprintf("http://%s/update/", cfg.ServerAddress)

		metrics.mu.RLock()
		arr := metrics.Arr
		metrics.mu.RUnlock()

		for _, metric := range arr {
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

			if cfg.HashKey != "" {
				hash := utils.GenerateHash(buf.Bytes(), cfg.HashKey)
				client.SetHeader("HashSHA256", hash)
			}

			requests <- ReportTask{
				Request: client.R().SetBody(buf.Bytes()),
				URL:     url,
			}
		}
		log.Logger.Info().Msg("finished reporting metrics to server/v2...")
	}
}

func reportMetricsV3(log logger.Logger, cfg configs.Config, metrics *Metrics, requests chan<- ReportTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
		log.Logger.Info().Msg("started reporting metrics to server/v3...")
		if len(metrics.Arr) == 0 {
			log.Logger.Info().Msg("no metrics to report")
			return
		}

		client := resty.New().
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Content-Type", "application/json")

		url := fmt.Sprintf("http://%s/updates/", cfg.ServerAddress)

		metrics.mu.RLock()
		arr := metrics.Arr
		metrics.mu.RUnlock()

		list := entity.MetricsList(arr)
		body, err := easyjson.Marshal(&list)
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot unmarshal metric object")
			return
		}

		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		_, err = gzipWriter.Write(body)
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot compress body")
			return
		}

		err = gzipWriter.Close()
		if err != nil {
			log.Logger.Info().Err(err).Msg("cannot close gzip writer")
			return
		}

		if cfg.HashKey != "" {
			hash := utils.GenerateHash(buf.Bytes(), cfg.HashKey)
			client.SetHeader("HashSHA256", hash)
		}

		requests <- ReportTask{
			Request: client.R().SetBody(buf.Bytes()),
			URL:     url,
		}

		log.Logger.Info().Msg("finished reporting metrics to server/v3...")
	}
}
