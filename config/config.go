package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ApiKey string `json:"apiKey"`
}

func LoadConfig(configFilepath string) (Config, error) {
	configJson, err := os.ReadFile(configFilepath)
	if err != nil {
		return Config{}, err
	}

	var userConfig Config
	if err := json.Unmarshal(configJson, &userConfig); err != nil {
		return Config{}, err
	}
	return userConfig, nil
}