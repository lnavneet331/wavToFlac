package config

import (
	"os"
)

type Config struct {
	ServerPort string
	LogLevel   string
}

func New() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", ":8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}