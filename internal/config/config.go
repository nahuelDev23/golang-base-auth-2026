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
	SessionDurationHours int
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	retention := os.Getenv("SESSION_RETENTION_DAYS")

	if retention == "" {
		retention = "90"
	}

	sessionRetentionDays, err := strconv.Atoi(retention)
	if err != nil {
		return nil, err
	}

	sessionDurationHours := os.Getenv("SESSION_DURATION_HOURS")

	if sessionDurationHours == "" {
		sessionDurationHours = "24"
	}

	sessionDuration, err := strconv.Atoi(sessionDurationHours)
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		Port:                 os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		SessionRetentionDays: sessionRetentionDays,
		SessionDurationHours: sessionDuration,
	}, nil
}
