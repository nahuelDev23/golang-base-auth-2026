package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL          string
	Port                 string
	JWTSecret            string
	SessionRetentionDays int
}

func Load() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	retention := os.Getenv("SESSION_RETENTION_DAYS")

	if retention == "" {
		retention = "90"
	}

	sessionRetentionDays, err := strconv.Atoi(retention)
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		Port:                 os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		SessionRetentionDays: sessionRetentionDays,
	}, nil
}
