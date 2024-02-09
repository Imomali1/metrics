package tasks

import (
	"bufio"
	"context"
	store "github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/mailru/easyjson"
	"os"
	"time"
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

func WriteMetricsToFile(ctx context.Context, storage *store.Storage, filename string, interval int) error {
	fw, err := NewFileWriter(filename)
	if err != nil {
		return err
	}

	storeTicker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-storeTicker.C:
			err = fw.WriteAllMetrics(storage)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (fw *FileWriter) WriteAllMetrics(storage *store.Storage) error {
	metrics, err := storage.Memory.ListMetrics()
	if err != nil {
		return err
	}

	var data []byte
	data, err = easyjson.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = fw.writer.Write(data)
	if err != nil {
		return err
	}

	return fw.writer.Flush()
}
