package utils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleGenerateHash() {
	out := GenerateHash([]byte("hello"), "world")
	fmt.Println(out)

	// Output:
	// 936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af
}

func TestGenerateHash(t *testing.T) {
	testCases := []struct {
		name     string
		value    []byte
		key      string
		expected string
	}{
		{
			name:     "Simple case",
			value:    []byte("hello"),
			key:      "world",
			expected: "936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af",
		},
		{
			name:     "Empty value",
			value:    []byte(""),
			key:      "key",
			expected: "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683",
		},
		{
			name:     "Empty key",
			value:    []byte("data"),
			key:      "",
			expected: "3a6eb0790f39ac87c94f3856b2dd2c5d110e6811602261a9a923d3bb23adc8b7",
		},
		{
			name:     "Both empty",
			value:    []byte(""),
			key:      "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GenerateHash(tc.value, tc.key)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func BenchmarkGenerateHash(b *testing.B) {
	testCases := []struct {
		value []byte
		key   string
	}{
		{
			value: []byte("small"),
			key:   "small",
		},
		{
			value: bytes.Repeat([]byte("medium"), 1024),
			key:   "medium",
		},
		{
			value: bytes.Repeat([]byte("large"), 1024*50),
			key:   "large",
		},
	}

	for _, tc := range testCases {
		b.Run(fmt.Sprintf("value size-%d bytes", len(tc.value)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = GenerateHash(tc.value, tc.key)
			}
		})
	}
}
