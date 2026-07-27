package main

import (
	"cobra-cli/cmd"
	"cobra-cli/config"
	"cobra-cli/storage/sqlite"
	"cobra-cli/task"
	"fmt"
	"log"
)

func main() {

	config := config.Load()

	store, err := sqlite.New(config.DB_PATH)

	if err != nil {
		log.Fatal("Error al iniciar la base de datos")
	}

	service := task.NewService(store)

	fmt.Println(service)

	cmd.Execute(service)
}
