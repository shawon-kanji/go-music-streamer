package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type PostgresDBConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	SSLMode  string `validate:"required"`
}

type Config struct {
	Postgres  PostgresDBConfig
	DBSources []string
	AppPort   string
	ENV       string
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Postgres: PostgresDBConfig{
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvWithDefault("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
		},
		DBSources: []string{"postgres"}, // Add more sources here as needed
		AppPort:   getEnvWithDefault("APP_PORT", "8080"),
		ENV:       getEnvWithDefault("ENV", "development"),
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
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
