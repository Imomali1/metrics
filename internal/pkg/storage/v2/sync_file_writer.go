package v2

import (
	"bufio"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func newFileWriter(filename string) (*fileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (f *fileWriter) Write(batch entity.MetricsList) error {
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
