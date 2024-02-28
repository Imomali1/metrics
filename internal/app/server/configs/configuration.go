package configs

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
	DatabaseDSN     string
	HashKey         string

	ServiceName string
	LogLevel    string
}

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true
	defaultDSN             = ""

	defaultServiceName = "metrics_server"
	defaultLogLevel    = "info"
)

func Parse(cfg *Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	storeInterval := flag.Int("i", defaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "полное имя файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", defaultRestore, "булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	databaseDSN := flag.String("d", defaultDSN, "адрес подключения к БД")
	hashKey := flag.String("k", "", "Ключ для подписи данных")

	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", serverAddress)
	cfg.StoreInterval = getEnvInt("STORE_INTERVAL", storeInterval)
	cfg.FileStoragePath = getEnvString("FILE_STORAGE_PATH", fileStoragePath)
	cfg.Restore = getEnvBool("RESTORE", restore)
	cfg.DatabaseDSN = getEnvString("DATABASE_DSN", databaseDSN)
	cfg.HashKey = getEnvString("KEY", hashKey)

	cfg.ServiceName = defaultServiceName
	cfg.LogLevel = defaultLogLevel
}

func getEnvString(key string, argumentValue *string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return *argumentValue
}

func getEnvInt(key string, argumentValue *int) int {
	envValue, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return envValue
	}
	return *argumentValue
}

func getEnvBool(key string, argumentValue *bool) bool {
	envValue, err := strconv.ParseBool(os.Getenv(key))
	if err == nil {
		return envValue
	}
	return *argumentValue
}
