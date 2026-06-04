package config

import (
	"os"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedpandaBrokers string
}

func Load() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8101"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "ahs"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedpandaBrokers: getEnv("REDPANDA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
