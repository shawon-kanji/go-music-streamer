package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSLMODE  string
	APP_PORT    string
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DB_HOST:     getEnvWithDefault("DB_HOST", "localhost"),
		DB_PORT:     getEnvWithDefault("DB_PORT", "5432"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_SSLMODE:  getEnvWithDefault("DB_SSLMODE", "disable"),
		APP_PORT:    getEnvWithDefault("APP_PORT", "8080"),
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	if cfg.DB_USER == "" {
		return Config{}, fmt.Errorf("missing DB_USER")
	}
	if cfg.DB_PASSWORD == "" {
		return Config{}, fmt.Errorf("missing DB_PASSWORD")
	}
	if cfg.DB_NAME == "" {
		return Config{}, fmt.Errorf("missing DB_NAME")
	}

	return cfg, nil
}

func getEnvWithDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
