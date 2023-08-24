package config

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		DB     `yaml:"db"`
		Server `yaml:"server"`
		Logger json.RawMessage
	}
	DB struct {
		DSN string `yaml:"dsn"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
)

func New() (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}
	err = cfg.loadJSON("./config/logger.json")
	if err != nil {
		return Config{}, fmt.Errorf("cfg.loadJSON: %w", err)
	}
	return cfg, nil
}

func (cfg *Config) loadJSON(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	if err = json.Unmarshal(bytes, cfg); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	return nil
}
