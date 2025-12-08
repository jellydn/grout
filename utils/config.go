package utils

import (
	"encoding/json"
	"fmt"
	"grout/models"
	"os"
	"time"

	gaba "github.com/UncleJunVIP/gabagool/v2/pkg/gabagool"
)

func LoadConfig() (*models.Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("reading config.json: %w", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config.json: %w", err)
	}

	if config.ApiTimeout == 0 {
		config.ApiTimeout = 30 * time.Minute
	}

	if config.DownloadTimeout == 0 {
		config.DownloadTimeout = 60 * time.Minute
	}

	return &config, nil
}

func SaveConfig(config *models.Config) error {
	if config.LogLevel == "" {
		config.LogLevel = "ERROR"
	}

	gaba.SetRawLogLevel(config.LogLevel)

	pretty, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		gaba.GetLogger().Error("Failed to marshal config to JSON", "error", err)
		return err
	}

	if err := os.WriteFile("config.json", pretty, 0644); err != nil {
		gaba.GetLogger().Error("Failed to write config file", "error", err)
		return err
	}

	return nil
}
