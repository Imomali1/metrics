package agent

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerAddress  string
	PollInterval   int
	ReportInterval int
}

func Parse(cfg *Config) {
	serverAddress := flag.String("a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", 2, "частота опроса метрик из пакета runtime")
	reportInterval := flag.Int("r", 10, "частота отправки метрик на сервер")
	flag.Parse()

	cfg.ServerAddress = castToString(getEnvOrDefaultValue("ADDRESS", *serverAddress))
	cfg.PollInterval = castToInt(getEnvOrDefaultValue("POLL_INTERVAL", *pollInterval))
	cfg.ReportInterval = castToInt(getEnvOrDefaultValue("REPORT_INTERVAL", *reportInterval))
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

func castToInt(i interface{}) int {
	v, ok := i.(int)
	if !ok {
		return 0
	}
	return v
}
