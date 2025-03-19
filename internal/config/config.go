package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (*Config, error) {
	var cfg Config

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		log.Fatal("HTTP_ADDR не установлен")
	}
	cfg.HTTPAddr = httpAddr

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}
	cfg.DatabaseURL = dbURL

	return &cfg, nil
}
