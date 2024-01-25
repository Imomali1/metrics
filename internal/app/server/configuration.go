package server

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	ServiceName   string
	LogLevel      string
}

const (
	defaultServerAddress = "localhost:8080"
	defaultServiceName   = "metrics_server"
	defaultLogLevel      = "info"
)

func Parse(cfg *Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", *serverAddress, defaultServerAddress)
	cfg.ServiceName = defaultServiceName
	cfg.LogLevel = defaultLogLevel
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
