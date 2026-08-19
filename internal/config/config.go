package config

import (
	"os"

	"github.com/AVZotov/draft-survey/internal/logger"
)

type Config struct {
	Port     string
	DBPath   string
	LogLevel logger.Level
	Version  string
}

func Load() Config {
	return Config{
		Port:     getEnv("PORT", ":3399"),
		DBPath:   getEnv("DB_PATH", "./data/draft-survey.db"),
		LogLevel: parseLogLevel(getEnv("LOG_LEVEL", "info")),
		Version:  getEnv("VERSION", "dev"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseLogLevel(s string) logger.Level {
	switch s {
	case "debug":
		return logger.LevelDebug
	case "warn":
		return logger.LevelWarn
	case "error":
		return logger.LevelError
	default:
		return logger.LevelInfo
	}
}
