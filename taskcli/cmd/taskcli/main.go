package main

import (
	"log"
	"taskcli/internal/cli"
	"taskcli/internal/config"
	"taskcli/internal/storage/sqlite"
	"taskcli/internal/task"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Println(err)
		log.Fatal("Error al cargar variables de entorno")
	}

	store, err := sqlite.New(config.DBPath)

	if err != nil {
		log.Println(err)
		log.Fatal("Error al establecer coneccion a la base de datos")
	}

	defer store.Close()

	taskService := task.NewService(store)

	cli.NewRootCommand(store, taskService)
}
