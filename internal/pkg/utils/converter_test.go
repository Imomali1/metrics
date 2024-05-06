package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "Float64",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "Integer",
			input:    42,
			expected: 42,
		},
		{
			name:     "Boolean",
			input:    true,
			expected: true,
		},
		{
			name:     "String",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "NilPointer",
			input:    (*string)(nil),
			expected: (*string)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ptr(tt.input)
			switch v := tt.expected.(type) {
			case float64:
				if res, ok := (*result).(float64); ok {
					assert.Equal(t, v, res, "Failed %s", tt.name)
				} else {
					t.Errorf("Type assertion to float64 failed")
				}
			case int:
				if res, ok := (*result).(int); ok {
					assert.Equal(t, v, res, "Failed %s", tt.name)
				} else {
					t.Errorf("Type assertion to int failed")
				}
			case bool:
				if res, ok := (*result).(bool); ok {
					assert.Equal(t, v, res, "Failed %s", tt.name)
				} else {
					t.Errorf("Type assertion to bool failed")
				}
			case string:
				if res, ok := (*result).(string); ok {
					assert.Equal(t, v, res, "Failed %s", tt.name)
				} else {
					t.Errorf("Type assertion to string failed")
				}
			case *string:
				if res, ok := (*result).(*string); ok {
					assert.Nil(t, res, "Failed %s, expected nil pointer", tt.name)
				} else {
					t.Errorf("Type assertion to string failed")
				}
			}
		})
	}
}
