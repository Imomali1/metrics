package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/cipher"
	"github.com/Imomali1/metrics/internal/pkg/utils"
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

func (a agent) reportMetricsV1(
	metrics *Metrics,
	requests chan<- ReportTask,
	wg *sync.WaitGroup,
	quit <-chan os.Signal,
) {
	ticker := time.NewTicker(a.cfg.ReportInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()

				a.log.Info().Msg("started reporting metrics to server/v1...")
				if len(metrics.Arr) == 0 {
					a.log.Info().Msg("no metrics to report")
					return
				}

				client := resty.New().SetHeader("Content-Type", "text/plain")

				metrics.mu.RLock()
				arr := metrics.Arr
				metrics.mu.RUnlock()

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

					requests <- ReportTask{
						Request: client.R(),
						URL:     url,
					}
				}
				a.log.Info().Msg("finished reporting metrics to server/v1...")
			}()
		case <-quit:
			a.log.Info().Msg("stopped reporting metrics to server/v1...")
			return
		}
	}
}

func (a agent) worker(tasks <-chan ReportTask) {
	for task := range tasks {
		if err := task.Process(); err != nil {
			a.log.Info().Err(err).Msg("error in reporting metrics to server")
		} else {
			a.log.Info().Msg("metrics reported successfully")
		}
	}
}

func (a agent) reportMetricsV2(
	metrics *Metrics,
	requests chan<- ReportTask,
	wg *sync.WaitGroup,
	quit <-chan os.Signal,
) {
	ticker := time.NewTicker(a.cfg.ReportInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()

				a.log.Info().Msg("started reporting metrics to server/v2...")
				if len(metrics.Arr) == 0 {
					a.log.Info().Msg("no metrics to report")
					return
				}

				client := resty.New().
					SetHeader("Content-Encoding", "gzip").
					SetHeader("Content-Type", "application/json")

				url := fmt.Sprintf("http://%s/update/", a.cfg.ServerAddress)

				metrics.mu.RLock()
				arr := metrics.Arr
				metrics.mu.RUnlock()

				for _, metric := range arr {
					body, err := easyjson.Marshal(metric)
					if err != nil {
						a.log.Info().Err(err).Msg("cannot unmarshal metric object")
						continue
					}

					if a.publicKey != nil {
						body, err = cipher.EncryptRSA(a.publicKey, body)
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

					requests <- ReportTask{
						Request: client.R().SetBody(buf.Bytes()),
						URL:     url,
					}
				}
				a.log.Info().Msg("finished reporting metrics to server/v2...")
			}()
		case <-quit:
			a.log.Info().Msg("stopped reporting metrics to server/v2...")
			return
		}
	}
}

func (a agent) reportMetricsV3(
	metrics *Metrics,
	requests chan<- ReportTask,
	wg *sync.WaitGroup,
	quit <-chan os.Signal,
) {
	ticker := time.NewTicker(a.cfg.ReportInterval * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()

				a.log.Info().Msg("started reporting metrics to server/v3...")
				if len(metrics.Arr) == 0 {
					a.log.Info().Msg("no metrics to report")
					return
				}

				client := resty.New().
					SetHeader("Content-Encoding", "gzip").
					SetHeader("Content-Type", "application/json")

				url := fmt.Sprintf("http://%s/updates/", a.cfg.ServerAddress)

				metrics.mu.RLock()
				arr := metrics.Arr
				metrics.mu.RUnlock()

				list := entity.MetricsList(arr)
				body, err := easyjson.Marshal(&list)
				if err != nil {
					a.log.Info().Err(err).Msg("cannot unmarshal metric object")
					return
				}

				if a.publicKey != nil {
					body, err = cipher.EncryptRSA(a.publicKey, body)
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

				requests <- ReportTask{
					Request: client.R().SetBody(buf.Bytes()),
					URL:     url,
				}

				a.log.Info().Msg("finished reporting metrics to server/v3...")
			}()
		case <-quit:
			a.log.Info().Msg("stopped reporting metrics to server/v3...")
			return
		}
	}
}
