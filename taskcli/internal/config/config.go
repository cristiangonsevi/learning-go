package config

import "os"

type Config struct {
	DBPath string
}

func Load() (*Config, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "taskcli.db"
	}
	return &Config{DBPath: path}, nil
}
