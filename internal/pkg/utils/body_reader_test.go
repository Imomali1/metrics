package utils

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadAll(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		maxSize  int
	}{
		{
			name:     "normal small input",
			input:    "hello, world",
			expected: "hello, world",
			maxSize:  1024,
		},
		{
			name:     "exact buffer size",
			input:    strings.Repeat("a", 10*1024),
			expected: strings.Repeat("a", 10*1024),
			maxSize:  10 * 1024,
		},
		{
			name:     "need to expand buffer",
			input:    strings.Repeat("a", 11*1024),
			expected: strings.Repeat("a", 11*1024),
			maxSize:  10 * 1024,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			maxSize:  1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ReadAll(reader)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

type MockReader struct {
	Data      []byte
	ReadDelay int // delay in microseconds per read, to simulate slow I/O
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	if len(m.Data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, m.Data)
	m.Data = m.Data[n:]
	if m.ReadDelay > 0 {
		time.Sleep(time.Microsecond * time.Duration(m.ReadDelay))
	}
	return n, nil
}

func BenchmarkReadAll(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "1KB",
			size: 1024,
		},
		{
			name: "10KB",
			size: 10 * 1024,
		},
		{
			name: "100KB",
			size: 100 * 1024,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			data := bytes.Repeat([]byte("a"), tt.size)
			reader := &MockReader{Data: data}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ReadAll(reader)
				if err != nil {
					b.Fatal(err)
				}
				reader.Data = data
			}
		})
	}
}
