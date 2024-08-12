package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type FileConfig struct {
	ServerAddress   *string        `json:"address"`
	StoreInterval   *time.Duration `json:"store_interval"`
	FileStoragePath *string        `json:"store_file"`
	Restore         *bool          `json:"restore"`
	DatabaseDSN     *string        `json:"database_dsn"`
	PrivateKeyPath  *string        `json:"crypto_key"`
}

func LoadFileConfig(configPath string) (FileConfig, error) {
	conf := FileConfig{}

	if configPath == "" {
		return conf, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return conf, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to parse config file: %w", err)
	}

	return conf, nil
}
