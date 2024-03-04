package utils

import "time"

var (
	maxAttempts = 3
	delays      = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
)

func DoWithRetries(fn func() error) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if err = fn(); err != nil {
			time.Sleep(delays[i])
			continue
		}
		return nil
	}
	return err
}
