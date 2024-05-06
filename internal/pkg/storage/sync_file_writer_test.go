package storage

import (
	"bufio"
	"os"
	"testing"

	"github.com/Imomali1/metrics/internal/entity"
)

func Test_fileWriter_Write(t *testing.T) {
	type fields struct {
		file   *os.File
		writer *bufio.Writer
	}
	type args struct {
		batch entity.MetricsList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileWriter{
				file:   tt.fields.file,
				writer: tt.fields.writer,
			}
			if err := f.Write(tt.args.batch); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newFileWriter(t *testing.T) {

}
