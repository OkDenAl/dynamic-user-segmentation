package config

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		DB                    `yaml:"db"`
		Server                `yaml:"server"`
		GDriveCredentialsPath string `env-required:"true" env:"GOOGLE_DRIVE_JSON_FILE_PATH"`
		Logger                json.RawMessage
	}
	DB struct {
		DSN string `env-required:"true" env:"PG_DSN"`
	}
	Server struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env - %w", err)
	}
	err = cfg.loadJSON("./config/logger.json")
	if err != nil {
		return nil, fmt.Errorf("error load json - %w", err)
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
