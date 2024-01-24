package logger

import (
	"github.com/rs/zerolog"
	"io"
)

type Logger struct {
	Logger zerolog.Logger
}

var log *Logger

func initLogger(writer io.Writer, level string) zerolog.Logger {
	// Parse Logger Level
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	// Set Global Logger Level
	zerolog.SetGlobalLevel(lvl)

	// Create Logger Instance
	l := zerolog.New(writer).With().Timestamp().Logger()
	log = &Logger{
		Logger: l,
	}

	return l
}

func NewLogger(writer io.Writer, level string, serviceName string) Logger {
	var l zerolog.Logger
	if log != nil {
		l = log.Logger
	} else {
		l = initLogger(writer, level)
	}
	return Logger{
		Logger: l.With().Str("service", serviceName).Logger(),
	}
}
