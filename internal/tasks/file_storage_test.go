package tasks

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/entity"
	"github.com/Imomali1/metrics/internal/pkg/storage"
	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func TestNewFileWriter(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		valid    bool
	}{
		{
			name:     "test-1",
			filename: "/tmp/test-1.json",
			valid:    true,
		},
		{
			name:     "test-2",
			filename: "/invalid/path/test-1.json",
			valid:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check valid case
			fw, err := NewFileWriter(tt.filename)
			if !tt.valid {
				require.Error(t, err)
				require.Nil(t, fw)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, fw)
			require.NotNil(t, fw.file)
			require.NotNil(t, fw.file)
		})
	}
}

func TestWriteMetricsToFile(t *testing.T) {
	metrics := entity.MetricsList{
		{ID: "Gauge1", MType: entity.Gauge, Value: utils.Ptr(123.0)},
		{ID: "Counter1", MType: entity.Counter, Delta: utils.Ptr(int64(123))},
	}

	tests := []struct {
		name    string
		store   storage.Storage
		wantErr bool
	}{
		{
			name:    "valid",
			store:   NewMockStorage(WithMetrics(metrics)),
			wantErr: false,
		},
		{
			name:    "invalid",
			store:   NewMockStorage(WithError()),
			wantErr: true,
		},
	}

	// Initialize mock FileWriter
	var b bytes.Buffer
	fw := &FileWriter{
		writer: bufio.NewWriter(&b),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storedMetrics, err := tt.store.GetAll(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			err = fw.WriteAllMetrics(storedMetrics)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expectedData := ``
			for _, metric := range metrics {
				data, _ := easyjson.Marshal(metric)
				expectedData += string(data) + "\n"
			}
			require.Equal(t, expectedData, b.String())
		})
	}
}
