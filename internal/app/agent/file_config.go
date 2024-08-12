package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type FileConfig struct {
	ServerAddress  *string        `json:"address"`
	PollInterval   *time.Duration `json:"poll_interval"`
	ReportInterval *time.Duration `json:"report_interval"`
	PublicKeyPath  *string        `json:"crypto_key"`
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
