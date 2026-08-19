package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost string
	DBName string
	DBUser string
	DBPass string
	DBPort int
}

func Load() (*Config, error) {

	config := &Config{
		DBHost: getEnvStr("DB_HOST", "localhost"),
		DBName: getEnvStr("DB_NAME", ""),
		DBUser: getEnvStr("DB_USER", ""),
		DBPass: getEnvStr("DB_PASS", ""),
		DBPort: getEnvInt("DB_PORT", 5432),
	}

	err := validate(config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func getEnvStr(key string, defaultValue string) string {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return v
	}
	return defaultValue
}

func validate(cfg *Config) error {
	errors := []string{}
	if cfg.DBHost == "" {
		errors = append(errors, "DB_HOST is required")
	}
	if cfg.DBName == "" {
		errors = append(errors, "DB_NAME is required")
	}
	if cfg.DBUser == "" {
		errors = append(errors, "DB_USER is required")
	}
	if cfg.DBPass == "" {
		errors = append(errors, "DB_PASS is required")
	}
	if cfg.DBPort <= 2000 || cfg.DBPort >= 65555 {
		errors = append(errors, "DB_PORT is out of range")
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("[ERROR] %v", strings.Join(errors, ","))
}
