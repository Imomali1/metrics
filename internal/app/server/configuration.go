package server

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress   string
	StoreInterval   int
	FileStoragePath string
	Restore         bool

	ServiceName string
	LogLevel    string
}

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true

	defaultServiceName = "metrics_server"
	defaultLogLevel    = "info"
)

func Parse(cfg *Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	storeInterval := flag.Int("i", defaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "полное имя файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", defaultRestore, "булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")

	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", serverAddress, defaultServerAddress)
	cfg.StoreInterval = getEnvInt("STORE_INTERVAL", storeInterval, defaultStoreInterval)
	cfg.FileStoragePath = getEnvString("FILE_STORAGE_PATH", fileStoragePath, defaultFileStoragePath)
	cfg.Restore = getEnvBool("RESTORE", restore, defaultRestore)

	cfg.ServiceName = defaultServiceName
	cfg.LogLevel = defaultLogLevel
}

func getEnvString(key string, argumentValue *string, defaultValue string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	if argumentValue != nil {
		return *argumentValue
	}
	return defaultValue
}

func getEnvInt(key string, argumentValue *int, defaultValue int) int {
	envValue, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return envValue
	}
	if argumentValue != nil {
		return *argumentValue
	}
	return defaultValue
}

func getEnvBool(key string, argumentValue *bool, defaultValue bool) bool {
	envValue, err := strconv.ParseBool(os.Getenv(key))
	if err == nil {
		return envValue
	}
	if argumentValue != nil {
		return *argumentValue
	}
	return defaultValue
}
