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
	if value, err := getEnvAsInt("DB_MAX_OPEN_CONNS"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.DBMaxOpenConns = value
	}
	if value, err := getEnvAsInt("DB_MAX_IDLE_CONNS"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.DBMaxIdleConns = value
	}
	if value, err := getEnvAsDuration("DB_CONN_MAX_LIFETIME"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.DBConnMaxLifetime = value
	}
	if value, err := getEnvAsDuration("DB_PING_TIMEOUT"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.DBPingTimeout = value
	}

	// Server timeouts
	if value, err := getEnvAsDuration("SERVER_READ_TIMEOUT"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.ServerReadTimeout = value
	}
	if value, err := getEnvAsDuration("SERVER_WRITE_TIMEOUT"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.ServerWriteTimeout = value
	}

	// CORS configuration
	cfg.CORSAllowOrigins = getEnvAsSlice("CORS_ALLOW_ORIGINS")
	cfg.CORSAllowMethods = getEnvAsSlice("CORS_ALLOW_METHODS")
	cfg.CORSAllowHeaders = getEnvAsSlice("CORS_ALLOW_HEADERS")
	if value, err := getEnvAsBool("CORS_ALLOW_CREDENTIALS"); err != nil {
		errors = append(errors, err.Error())
	} else {
		cfg.CORSAllowCredentials = value
	}

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

func getEnvAsInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to 0", key)
		return 0, nil
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value for %s: %w", key, err)
	}
	return result, nil
}

func getEnvAsDuration(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to 0", key)
		return 0, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value for %s: %w", key, err)
	}
	return duration, nil
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

func getEnvAsBool(key string) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("warning: %s not set, defaulting to false", key)
		return false, nil
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid boolean value for %s: %w", key, err)
	}
	return result, nil
}
