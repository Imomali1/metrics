package agent

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress  string
	PollInterval   int
	ReportInterval int
}

const (
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

func Parse(cfg *Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", defaultPollInterval, "частота опроса метрик из пакета runtime")
	reportInterval := flag.Int("r", defaultReportInterval, "частота отправки метрик на сервер")
	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", *serverAddress, defaultServerAddress)
	cfg.PollInterval = getEnvInt("POLL_INTERVAL", *pollInterval, defaultPollInterval)
	cfg.ReportInterval = getEnvInt("REPORT_INTERVAL", *reportInterval, defaultReportInterval)
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

func getEnvInt(key string, argumentValue int, defaultValue int) int {
	if os.Getenv(key) != "" {
		value, err := strconv.Atoi(os.Getenv(key))
		if err == nil {
			return value
		}
	}

	if argumentValue > 0 {
		return argumentValue
	}
	return defaultValue
}
