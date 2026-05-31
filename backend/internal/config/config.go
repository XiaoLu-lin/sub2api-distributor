package config

import (
	"fmt"
	"os"
	"strings"
)

// Config contains all runtime settings needed by the distributor backend.
type Config struct {
	ServerPort         string
	DatabaseDSN        string
	JWTSecret          string
	AppEnv             string
	CORSAllowedOrigins []string
	StaticDir          string
}

// Load reads runtime configuration from environment variables and applies
// development defaults when no override is provided.
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:         getEnv("SERVER_PORT", "8091"),
		DatabaseDSN:        getEnv("DATABASE_DSN", "postgres://sub2api@localhost:5432/sub2api?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "sub2api-distributor-dev-secret"),
		AppEnv:             getEnv("APP_ENV", "development"),
		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "")),
		StaticDir:          strings.TrimSpace(getEnv("STATIC_DIR", "")),
	}

	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

// getEnv returns the configured environment value or the provided fallback.
func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
