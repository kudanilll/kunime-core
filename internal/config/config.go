package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	Env              string
	ClientAPIKey     string
	AdminAPIKey      string
	KunimeAPIBaseURL string
	KunimeAPIKey     string
}

var cfg *Config

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg = &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		Env:              getEnv("APP_ENV", "development"),
		ClientAPIKey:     getEnv("CLIENT_API_KEY", ""),
		AdminAPIKey:      getEnv("ADMIN_API_KEY", ""),
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

func ClientAPIKey() string {
	if cfg == nil {
		return ""
	}
	return cfg.ClientAPIKey
}

func AdminAPIKey() string {
	if cfg == nil {
		return ""
	}
	return cfg.AdminAPIKey
}
