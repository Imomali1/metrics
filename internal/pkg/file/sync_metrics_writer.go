package file

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type SyncFileWriter interface {
	Write(batch entity.MetricsList) error
}

type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewSyncMetricsWriter(filename string) (SyncFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (f *fileWriter) Write(batch entity.MetricsList) error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	for _, one := range batch {
		data, err := easyjson.Marshal(one)
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
