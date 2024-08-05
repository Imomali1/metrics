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
	HashKey        string
	RateLimit      int
	PublicKeyPath  string

	LogLevel    string
	ServiceName string
}

const (
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultLogLevel       = "info"
	defaultServiceName    = "metrics_agent"
)

func LoadConfig() (cfg Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", defaultPollInterval, "частота опроса метрик из пакета runtime")
	reportInterval := flag.Int("r", defaultReportInterval, "частота отправки метрик на сервер")
	hashKey := flag.String("k", "", "Ключ для подписи данных")
	rateLimit := flag.Int("l", 1, "количество одновременно исходящих запросов на сервер")
	publicKeyPath := flag.String("crypto-key", "", "путь до файла с публичным ключом")
	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", serverAddress)
	cfg.PollInterval = getEnvInt("POLL_INTERVAL", pollInterval)
	cfg.ReportInterval = getEnvInt("REPORT_INTERVAL", reportInterval)
	cfg.HashKey = getEnvString("KEY", hashKey)
	cfg.RateLimit = getEnvInt("RATE_LIMIT", rateLimit)
	cfg.PublicKeyPath = getEnvString("CRYPTO_KEY", publicKeyPath)

	cfg.LogLevel = defaultLogLevel
	cfg.ServiceName = defaultServiceName

	return cfg
}

func getEnvString(key string, argumentValue *string) string {
	envValue, exists := os.LookupEnv(key)
	if !exists {
		return *argumentValue
	}
	return envValue
}

func getEnvInt(key string, argumentValue *int) int {
	envValue, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return envValue
	}
	return *argumentValue
}
