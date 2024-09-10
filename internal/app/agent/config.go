package agent

import (
	"flag"
	"os"
	"strconv"

	"github.com/Imomali1/metrics/internal/pkg/utils"
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
	serverAddress := flag.String("a", "", "отвечает за адрес эндпоинта HTTP-сервера")
	pollInterval := flag.Int("p", 0, "частота опроса метрик из пакета runtime")
	reportInterval := flag.Int("r", 0, "частота отправки метрик на сервер")
	hashKey := flag.String("k", "", "Ключ для подписи данных")
	rateLimit := flag.Int("l", 1, "количество одновременно исходящих запросов на сервер")
	publicKeyPath := flag.String("crypto-key", "", "путь до файла с публичным ключом")
	shortConfigFilePath := flag.String("c", "", "путь до файла конфигурации short")
	longConfigFilePath := flag.String("config", "", "путь до файла конфигурации long")

	flag.Parse()

	var configFilePath string

	if *shortConfigFilePath != "" {
		configFilePath = *shortConfigFilePath
	}

	if *longConfigFilePath != "" {
		configFilePath = *longConfigFilePath
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

	var filePollInterval *int
	if fileConf.PollInterval != nil {
		filePollInterval = utils.Ptr(int(fileConf.PollInterval.Seconds()))
	}

	cfg.PollInterval = getEnvInt(
		"POLL_INTERVAL",
		*pollInterval,
		filePollInterval,
		defaultPollInterval,
	)

	var fileReportInterval *int
	if fileConf.ReportInterval != nil {
		fileReportInterval = utils.Ptr(int(fileConf.ReportInterval.Seconds()))
	}

	cfg.ReportInterval = getEnvInt(
		"REPORT_INTERVAL",
		*reportInterval,
		fileReportInterval,
		defaultReportInterval,
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

	if flagValue != "" {
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
	envValue, err := strconv.Atoi(os.Getenv(envKey))
	if err == nil {
		return envValue
	}

	if flagValue != 0 {
		return flagValue
	}

	if fileConfValue != nil {
		return *fileConfValue
	}

	return defaultValue
}
