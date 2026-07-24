package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"taskcli/internal/storage"
	"taskcli/internal/task"
)

func NewRootCommand(s storage.Store, taskService *task.Service) {
	subCommand := flag.NewFlagSet("subcommand", flag.ExitOnError)
	taskName := subCommand.String("name", "", "Nombre de la tarea")

	if len(os.Args) < 2 {
		fmt.Printf("Mostrando ayudame de cli\n")
		return
	}

	switch os.Args[1] {
	case "list":
		newListCommand(taskService)
	case "create":
		subCommand.Parse(os.Args[2:])
		createCommand(taskService, taskName)
	default:
		log.Printf("Comando `%s` no soportado\n", os.Args[1])
	}

}
