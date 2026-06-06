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

	MinIOEndpoint        string
	MinIOAccessKey       string
	MinIOSecretKey       string
	MinIOSnapshotsBucket string
	MinIOUseSSL          bool
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

		MinIOEndpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:       getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:       getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinIOSnapshotsBucket: getEnv("MINIO_SNAPSHOTS_BUCKET", "snapshots"),
		MinIOUseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
