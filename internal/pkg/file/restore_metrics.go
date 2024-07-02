package file

import (
	"bufio"
	"context"
	"errors"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/mailru/easyjson"
	"os"
)

func RestoreMetrics(ctx context.Context, filename string, store storage.Storage) error {
	if store == nil {
		return errors.New("storage not initialized")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var i int
	var metrics entity.MetricsList
	m := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Bytes()

		var metric entity.Metrics
		if err = easyjson.Unmarshal(line, &metric); err != nil {
			return err
		}

		idx, ok := m[metric.ID]
		if !ok {
			metrics = append(metrics, metric)
			m[metric.ID] = i
			i++
		} else {
			metrics[idx] = metric
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return store.Update(ctx, metrics)
}
