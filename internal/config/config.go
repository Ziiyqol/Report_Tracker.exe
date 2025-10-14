package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Theme string `json:"theme"` // "dark" или "light"
}

const configFile = "config.json"

func Load() Config {
	var cfg Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{Theme: "dark"}
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{Theme: "dark"}
	}
	return cfg
}

func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}
