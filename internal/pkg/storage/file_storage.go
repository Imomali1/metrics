package storage

import (
	"bufio"
	"github.com/mailru/easyjson"
	"os"
)

type FileStorage interface {
	WriteCounter(name string, delta int64) error
	WriteGauge(name string, value float64) error
	WriteAllMetrics() error
}

type fileStorage struct {
	metricStorage *Storage
	file          *os.File
	writer        *bufio.Writer
}

func newFileStorage(path string, metricStorage *Storage) (*fileStorage, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		metricStorage: metricStorage,
		file:          file,
		writer:        bufio.NewWriter(file),
	}, nil
}

func (f *fileStorage) WriteCounter(name string, delta int64) error {

}

func (f *fileStorage) WriteGauge(name string, value float64) error {

}

func (f *fileStorage) WriteAllMetrics() error {
	metrics, err := f.metricStorage.ListMetrics()
	if err != nil {
		return err
	}

	var data []byte
	data, err = easyjson.Marshal(metrics)
	if err != nil {
		return err
	}

	if _, err = f.writer.Write(data); err != nil {
		return err
	}

	return f.writer.Flush()
}
