package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	Env              string
	APIKey           string
	KunimeAPIBaseURL string
	KunimeAPIKey     string
}

var cfg *Config

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg = &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		Env:              getEnv("APP_ENV", "development"),
		APIKey:           getEnv("API_KEY", ""),
		KunimeAPIBaseURL: getEnv("KUNIME_API_BASE_URL", ""),
		KunimeAPIKey:     getEnv("KUNIME_API_KEY", ""),
	}

	return cfg, nil
}

func AppPort() string {
	if cfg == nil {
		return "8080"
	}
	return cfg.AppPort
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func APIKey() string {
	if cfg == nil {
		return ""
	}
	return cfg.APIKey
}
