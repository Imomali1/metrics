package logger

import (
	"github.com/rs/zerolog"
	"io"
)

func New(writer io.Writer, level string) (zerolog.Logger, error) {
	// Parse Log Level
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	// Set Global Log Level
	zerolog.SetGlobalLevel(lvl)

	zl := zerolog.New(writer).With().Timestamp().Logger()

	return zl, nil
}
