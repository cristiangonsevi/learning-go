package main

import (
	"cobra-cli/config"
	"cobra-cli/internal/cmd"
	"cobra-cli/internal/storage/sqlite"
	"cobra-cli/internal/task"
	"log"
)

func main() {

	config, err := config.Load()

	if err != nil {
		log.Fatal("Error al cargar variables de entorno", err)
	}

	store, errDb := sqlite.New(config.DB_PATH)

	if errDb != nil {
		log.Fatal("Error al iniciar la base de datos ", errDb)
	}

	defer store.Close()

	service := task.NewService(store)

	cmd.Execute(service)
}
