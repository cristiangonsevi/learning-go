package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PATH string
}

func Load() *Config {

	var envars map[string]string

	envars, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error al cargar variable de entorno")
	}

	return &Config{DB_PATH: envars["DB_PATH"]}
}
