package v1

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

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

func (f *fileStorage) WriteMetrics(metrics []entity.Metrics) error {
	for _, metric := range metrics {
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
	}

	return f.writer.Flush()
}
