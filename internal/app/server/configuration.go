package server

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	LogLevel      string
}

const (
	defaultServerAddress = "localhost:8080"
	defaultLogLevel      = "info"
)

func Parse(cfg *Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	logLevel := flag.String("log", defaultLogLevel, "устанавливает глобальный уровень логгера")
	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", *serverAddress, defaultServerAddress)
	cfg.LogLevel = getEnvString("LOG_LEVEL", *logLevel, defaultLogLevel)
}

func getEnvString(key string, argumentValue string, defaultValue string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	if argumentValue != "" {
		return argumentValue
	}
	return defaultValue
}
