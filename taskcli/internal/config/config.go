package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath string
}

func Load() (*Config, error) {

	var envars map[string]string

	envars, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error al leer archivo .env")
	}

	path := envars["DB_PATH"]

	log.Println("Environment ", path)
	if path == "" {
		log.Fatal("DB_PATH debe ser definido")
	}
	return &Config{DBPath: path}, nil
}
