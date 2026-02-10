package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg, nil
}
