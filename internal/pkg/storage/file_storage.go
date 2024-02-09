package storage

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type FileStorage interface {
	WriteCounter(name string, delta int64) error
	WriteGauge(name string, value float64) error
}

type fileStorage struct {
	file   *os.File
	writer *bufio.Writer
}

func newFileStorage(path string) (*fileStorage, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (f *fileStorage) WriteCounter(name string, delta int64) error {
	metric := entity.Metrics{
		ID:    name,
		MType: entity.Counter,
		Delta: &delta,
	}

	data, err := easyjson.Marshal(metric)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(data)
	if err != nil {
		return err
	}

	if err = f.writer.WriteByte('\n'); err != nil {
		return err
	}

	return f.writer.Flush()
}

func (f *fileStorage) WriteGauge(name string, value float64) error {
	metric := entity.Metrics{
		ID:    name,
		MType: entity.Gauge,
		Value: &value,
	}

	data, err := easyjson.Marshal(metric)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(data)
	if err != nil {
		return err
	}

	if err = f.writer.WriteByte('\n'); err != nil {
		return err
	}

	return f.writer.Flush()
}
