package agent

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func Parse(cfg *Config) {
	serverAddress := flag.String("a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", 2, "частота отправки метрик на сервер")
	reportInterval := flag.Int("r", 10, "частота опроса метрик из пакета runtime")
	flag.Parse()

	cfg.ServerAddress = castToString(getEnvOrDefaultValue("ADDRESS", *serverAddress))
	cfg.PollInterval = castToDuration(getEnvOrDefaultValue("POLL_INTERVAL", *pollInterval))
	cfg.ReportInterval = castToDuration(getEnvOrDefaultValue("REPORT_INTERVAL", *reportInterval))
}

func getEnvOrDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return os.Getenv(key)
}

func castToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func castToDuration(i interface{}) time.Duration {
	var t time.Duration
	switch v := i.(type) {
	case int:
		t = time.Duration(v) * time.Second
	case float64:
		t = time.Duration(v) * time.Second
	}
	return t
}
