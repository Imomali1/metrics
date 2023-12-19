package server

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
}

func Parse(cfg *Config) {
	serverAddress := flag.String("a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера")
	flag.Parse()

	cfg.ServerAddress = castToString(getEnvOrDefaultValue("ADDRESS", *serverAddress))
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
