package agent

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
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
	shortConfigFilePath := flag.String("c", "", "путь до файла конфигурации short")
	longConfigFilePath := flag.String("config", "", "путь до файла конфигурации long")

	flag.Parse()

	configFilePath := *longConfigFilePath

	if *shortConfigFilePath != "" {
		configFilePath = *shortConfigFilePath
	}

	path, found := os.LookupEnv("CONFIG")
	if found {
		configFilePath = path
	}

	fileConf, err := LoadFileConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	cfg.ServerAddress = getEnvString(
		"ADDRESS",
		*serverAddress,
		fileConf.ServerAddress,
		defaultServerAddress,
	)

	cfg.PollInterval = getEnvDuration(
		"POLL_INTERVAL",
		time.Duration(*pollInterval)*time.Second,
		fileConf.PollInterval,
		time.Duration(defaultPollInterval)*time.Second,
	)

	cfg.ReportInterval = getEnvDuration(
		"REPORT_INTERVAL",
		time.Duration(*reportInterval)*time.Second,
		fileConf.ReportInterval,
		time.Duration(defaultReportInterval)*time.Second,
	)

	cfg.HashKey = getEnvString(
		"KEY",
		*hashKey,
		nil,
		"",
	)

	cfg.RateLimit = getEnvInt(
		"RATE_LIMIT",
		*rateLimit,
		nil,
		1,
	)

	cfg.PublicKeyPath = getEnvString(
		"CRYPTO_KEY",
		*publicKeyPath,
		fileConf.PublicKeyPath,
		"",
	)

	cfg.LogLevel = defaultLogLevel
	cfg.ServiceName = defaultServiceName

	return cfg
}

func getEnvString(
	envKey string,
	flagValue string,
	fileConfValue *string,
	defaultValue string,
) string {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		return envValue
	}

	if flagValue != defaultValue {
		return flagValue
	}

	if fileConfValue != nil {
		return *fileConfValue
	}

	return defaultValue
}

func getEnvInt(
	envKey string,
	flagValue int,
	fileConfValue *int,
	defaultValue int,
) int {
	envValueStr, exists := os.LookupEnv(envKey)
	if exists {
		envValue, err := strconv.Atoi(envValueStr)
		if err == nil {
			return envValue
		}
	}

	if flagValue != defaultValue {
		return flagValue
	}

	if fileConfValue != nil {
		return *fileConfValue
	}

	return defaultValue
}

func getEnvDuration(
	key string,
	flagValue time.Duration,
	fileValue *time.Duration,
	defaultValue time.Duration,
) time.Duration {
	envValue, err := time.ParseDuration(os.Getenv(key))
	if err == nil {
		return envValue
	}

	if flagValue != defaultValue {
		return flagValue
	}

	if fileValue != nil {
		return *fileValue
	}

	return defaultValue
}
