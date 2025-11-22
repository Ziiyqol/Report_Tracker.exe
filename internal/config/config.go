package config

import (
	"encoding/json"
	"os"
)

const configFile = "config.json"

type Config struct {
	Theme string `json:"theme"` // "dark" или "light"
}

// Load загружает конфиг или возвращает дефолтный
func Load() Config {
	var cfg Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{Theme: "dark"} // Дефолтная тема
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{Theme: "dark"}
	}
	return cfg
}

// Save сохраняет конфиг
func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}
