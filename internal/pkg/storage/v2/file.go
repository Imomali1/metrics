package v2

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Imomali1/metrics/internal/entity"
	"github.com/mailru/easyjson"
	"os"
)

type reader struct {
	file    *os.File
	scanner *bufio.Scanner
}

type writer struct {
	file   *os.File
	writer *bufio.Writer
}

type fileStorage struct {
	r reader
	w writer
}

func newFileStorage(path string) (*fileStorage, error) {
	fileReader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	fileWriter, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		r: reader{
			file:    fileReader,
			scanner: bufio.NewScanner(fileReader),
		},
		w: writer{
			file:   fileWriter,
			writer: bufio.NewWriter(fileWriter),
		},
	}, nil
}

func (s *fileStorage) Update(ctx context.Context, batch entity.MetricsList) error {
	for _, one := range batch {
		data, err := easyjson.Marshal(one)
		if err != nil {
			return err
		}

		_, err = s.w.writer.Write(data)
		if err != nil {
			return err
		}

		if err = s.w.writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	return s.w.writer.Flush()
}

func (s *fileStorage) GetOne(ctx context.Context, id string, mType string) (entity.Metrics, error) {
	var wanted entity.Metrics
	for s.r.scanner.Scan() {
		line := s.r.scanner.Bytes()
		var metric entity.Metrics
		if err := easyjson.Unmarshal(line, &metric); err != nil {
			return entity.Metrics{}, err
		}
		if metric.ID == id && metric.MType == mType {
			wanted = entity.Metrics{
				ID:    id,
				MType: mType,
				Delta: metric.Delta,
				Value: metric.Value,
			}
			break
		}
	}

	if err := s.r.scanner.Err(); err != nil {
		return entity.Metrics{}, err
	}

	return wanted, nil
}

func (s *fileStorage) GetAll(ctx context.Context) (entity.MetricsList, error) {
	var metrics entity.MetricsList
	for s.r.scanner.Scan() {
		line := s.r.scanner.Bytes()
		var metric entity.Metrics
		if err := easyjson.Unmarshal(line, &metric); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	if err := s.r.scanner.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *fileStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("storage instance is not db, it is file based")
}

func (s *fileStorage) Close() {
	s.r.file.Close()
	s.w.file.Close()
}
