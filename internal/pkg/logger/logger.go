package logger

import (
	"github.com/rs/zerolog"
	"io"
)

type Logger struct {
	zerolog.Logger
}

var defaultLogger *Logger

func initLogger(writer io.Writer, level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	log := zerolog.New(writer).With().Timestamp().Logger()

	defaultLogger = &Logger{
		Logger: log,
	}

	return log
}

func NewLogger(writer io.Writer, level, service string) Logger {
	var log zerolog.Logger
	if defaultLogger != nil {
		log = defaultLogger.Logger
	} else {
		log = initLogger(writer, level)
	}
	return Logger{
		Logger: log.With().Str("service", service).Logger(),
	}
}
