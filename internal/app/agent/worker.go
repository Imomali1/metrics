package agent

import (
	"github.com/Imomali1/metrics/internal/pkg/utils"
	"github.com/go-resty/resty/v2"
)

type Job struct {
	Request *resty.Request
	URL     string
}

func (t *Job) Process() error {
	err := utils.DoWithRetries(func() error {
		_, err := t.Request.Post(t.URL)
		return err
	})
	return err
}

func (a *agent) worker() {
	for job := range a.jobsChan {
		if err := job.Process(); err != nil {
			a.log.Info().Err(err).Msg("error in reporting metrics to server")
		} else {
			a.log.Info().Msg("metrics reported successfully")
		}
	}
}
