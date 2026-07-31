package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	DB_PATH string
}

func Load() (*Config, error) {

	var envars map[string]string

	envars, err := godotenv.Read(".env")

	if err != nil {
		return nil, err
	}

	return &Config{DB_PATH: envars["DB_PATH"]}, nil
}
