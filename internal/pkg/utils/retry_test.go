package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoWithRetries(t *testing.T) {
	t.Run("Success on first try", func(t *testing.T) {
		var attempts int
		err := DoWithRetries(func() error {
			attempts++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("Success on third try", func(t *testing.T) {
		var attempts int
		err := DoWithRetries(func() error {
			attempts++
			if attempts < 3 {
				return errors.New("error")
			}
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("Fail after max attempts", func(t *testing.T) {
		var attempts int
		err := DoWithRetries(func() error {
			attempts++
			return errors.New("persistent error")
		})
		assert.Error(t, err)
		assert.Equal(t, maxAttempts, attempts)
	})
}
