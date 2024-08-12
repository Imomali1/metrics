package server

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/Imomali1/metrics/internal/api"
)

type Config struct {
	ServerAddress   string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
	API             api.Config
	PrivateKeyPath  string

	ServiceName string
	LogLevel    string
}

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300 * time.Second
	defaultFileStoragePath = "/tmp/metrics-database.json"
	defaultRestore         = true
	defaultDSN             = ""

	defaultServiceName = "metrics_server"
	defaultLogLevel    = "info"
)

func LoadConfig() (cfg Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	storeInterval := flag.Duration("i", defaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "полное имя файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", defaultRestore, "булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	databaseDSN := flag.String("d", defaultDSN, "адрес подключения к БД")
	hashKey := flag.String("k", "", "Ключ для подписи данных")
	privateKeyPath := flag.String("crypto-key", "", "путь до файла с приватным ключом")
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

	cfg.StoreInterval = getEnvDuration(
		"STORE_INTERVAL",
		*storeInterval,
		fileConf.StoreInterval,
		defaultStoreInterval,
	)

	cfg.FileStoragePath = getEnvString(
		"FILE_STORAGE_PATH",
		*fileStoragePath,
		fileConf.FileStoragePath,
		defaultFileStoragePath,
	)

	cfg.Restore = getEnvBool(
		"RESTORE",
		*restore,
		fileConf.Restore,
		defaultRestore,
	)

	cfg.DatabaseDSN = getEnvString(
		"DATABASE_DSN",
		*databaseDSN,
		fileConf.DatabaseDSN,
		defaultDSN,
	)

	cfg.PrivateKeyPath = getEnvString(
		"CRYPTO_KEY",
		*privateKeyPath,
		fileConf.PrivateKeyPath,
		"",
	)

	cfg.API.HashKey = getEnvString("KEY", *hashKey, nil, "")

	cfg.ServiceName = defaultServiceName
	cfg.LogLevel = defaultLogLevel

	return cfg
}

func getEnvString(
	key string,
	flagValue string,
	fileValue *string,
	defaultValue string,
) string {
	envValue, found := os.LookupEnv(key)
	if found {
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

func getEnvInt(
	key string,
	flagValue int,
	fileValue *int,
	defaultValue int,
) int {
	envValue, err := strconv.Atoi(os.Getenv(key))
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

func getEnvBool(
	key string,
	flagValue bool,
	fileValue *bool,
	defaultValue bool,
) bool {
	envValue, err := strconv.ParseBool(os.Getenv(key))
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
