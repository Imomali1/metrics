package tasks

import (
	"bufio"
	"context"
	"github.com/Imomali1/metrics/internal/entity"
	"os"
	"time"

	"github.com/mailru/easyjson"

	"github.com/Imomali1/metrics/internal/pkg/storage"
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func WriteMetricsToFile(
	ctx context.Context,
	store storage.Storage,
	filename string,
	interval int,
) error {
	fw, err := NewFileWriter(filename)
	if err != nil {
		return err
	}

	storeTicker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-storeTicker.C:
			var metrics entity.MetricsList
			metrics, err = store.GetAll(context.Background())
			if err != nil {
				return err
			}

			err = fw.WriteAllMetrics(metrics)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (fw *FileWriter) WriteAllMetrics(metrics entity.MetricsList) error {
	for _, metric := range metrics {
		data, err := easyjson.Marshal(metric)
		if err != nil {
			return err
		}

		_, err = fw.writer.Write(data)
		if err != nil {
			return err
		}

		if err = fw.writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	return fw.writer.Flush()
}
