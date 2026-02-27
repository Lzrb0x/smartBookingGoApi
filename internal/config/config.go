package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Application
	AppEnv string
	Port   string

	// Database
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBPingTimeout     time.Duration

	// Server
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration

	// CORS
	CORSAllowOrigins     []string
	CORSAllowMethods     []string
	CORSAllowHeaders     []string
	CORSAllowCredentials bool
}

func Load() (*Config, error) {
	cfg := &Config{}
	var errors []string

	// Variáveis obrigatórias
	cfg.AppEnv = mustGetEnv("APP_ENV", &errors)
	cfg.Port = mustGetEnv("PORT", &errors)
	cfg.DatabaseURL = mustGetEnv("DATABASE_URL", &errors)

	// Database configuration
	cfg.DBMaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS")
	cfg.DBMaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS")
	cfg.DBConnMaxLifetime = getEnvAsDuration("DB_CONN_MAX_LIFETIME")
	cfg.DBPingTimeout = getEnvAsDuration("DB_PING_TIMEOUT")

	// Server timeouts
	cfg.ServerReadTimeout = getEnvAsDuration("SERVER_READ_TIMEOUT")
	cfg.ServerWriteTimeout = getEnvAsDuration("SERVER_WRITE_TIMEOUT")

	// CORS configuration
	cfg.CORSAllowOrigins = getEnvAsSlice("CORS_ALLOW_ORIGINS")
	cfg.CORSAllowMethods = getEnvAsSlice("CORS_ALLOW_METHODS")
	cfg.CORSAllowHeaders = getEnvAsSlice("CORS_ALLOW_HEADERS")
	cfg.CORSAllowCredentials = getEnvAsBool("CORS_ALLOW_CREDENTIALS")

	if len(errors) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(errors, ", "))
	}

	return cfg, nil
}

func mustGetEnv(key string, errors *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*errors = append(*errors, key)
		return ""
	}
	return value
}

func getEnvAsInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to 0", key)
		return 0
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid integer value for %s: %v", key, err)
	}
	return result
}

func getEnvAsDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to 0", key)
		return 0
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid duration value for %s: %v", key, err)
	}
	return duration
}

func getEnvAsSlice(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to empty slice", key)
		return []string{}
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

func getEnvAsBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to false", key)
		return false
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("invalid boolean value for %s: %v", key, err)
	}
	return result
}
