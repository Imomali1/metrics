package agent

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileConfig struct {
	ServerAddress  *string `json:"address"`
	PollInterval   *int    `json:"poll_interval"`
	ReportInterval *int    `json:"report_interval"`
	PublicKeyPath  *string `json:"crypto_key"`
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
