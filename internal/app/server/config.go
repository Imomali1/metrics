package server

import (
	"flag"
	"os"
	"strconv"

	"github.com/Imomali1/metrics/internal/api"
)

type Config struct {
	ServerAddress   string
	StoreInterval   int
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
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-database.json"
	defaultRestore         = true
	defaultDSN             = ""

	defaultServiceName = "metrics_server"
	defaultLogLevel    = "info"
)

func LoadConfig() (cfg Config) {
	serverAddress := flag.String("a", defaultServerAddress, "отвечает за адрес эндпоинта HTTP-сервера")
	storeInterval := flag.Int("i", defaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "полное имя файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", defaultRestore, "булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	databaseDSN := flag.String("d", defaultDSN, "адрес подключения к БД")
	hashKey := flag.String("k", "", "Ключ для подписи данных")
	privateKeyPath := flag.String("crypto-key", "", "путь до файла с приватным ключом")
	flag.Parse()

	cfg.ServerAddress = getEnvString("ADDRESS", serverAddress)
	cfg.StoreInterval = getEnvInt("STORE_INTERVAL", storeInterval)
	cfg.FileStoragePath = getEnvString("FILE_STORAGE_PATH", fileStoragePath)
	cfg.Restore = getEnvBool("RESTORE", restore)
	cfg.DatabaseDSN = getEnvString("DATABASE_DSN", databaseDSN)
	cfg.API.HashKey = getEnvString("KEY", hashKey)
	cfg.PrivateKeyPath = getEnvString("CRYPTO_KEY", privateKeyPath)

	cfg.ServiceName = defaultServiceName
	cfg.LogLevel = defaultLogLevel

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

func getEnvBool(key string, argumentValue *bool) bool {
	envValue, err := strconv.ParseBool(os.Getenv(key))
	if err == nil {
		return envValue
	}
	return *argumentValue
}
