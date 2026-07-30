package main

import (
	"cobra-cli/config"
	"cobra-cli/internal/cmd"
	"cobra-cli/internal/storage/sqlite"
	"cobra-cli/internal/task"
	"log"
)

func main() {

	config := config.Load()

	store, err := sqlite.New(config.DB_PATH)

	if err != nil {
		log.Fatal("Error al iniciar la base de datos ", err)
	}

	service := task.NewService(store)

	cmd.Execute(service)
}
