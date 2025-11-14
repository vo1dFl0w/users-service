package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env"`
	HTTPaddr string `yaml:"http_addr"`
	DB       struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"db"`
	Secret string `yaml:"secret"`
}

// Load config from config.yaml
func LoadConfig() (*Config, error) {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return nil, fmt.Errorf("CONFIG_PATH not set")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
