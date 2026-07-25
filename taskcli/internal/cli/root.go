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
	createSubCommand := flag.NewFlagSet("createSubcommand", flag.ExitOnError)
	taskName := createSubCommand.String("name", "", "Nombre de la tarea")

	deleteSubCommand := flag.NewFlagSet("deleteSubCommand", flag.ExitOnError)
	taskId := deleteSubCommand.Int("id", 0, "Id de la tarea")

	if len(os.Args) < 2 {
		fmt.Printf("Mostrando ayudame de cli\n")
		return
	}

	switch os.Args[1] {
	case "list":
		newListCommand(taskService)
	case "create":
		createSubCommand.Parse(os.Args[2:])
		createCommand(taskService, taskName)
	case "delete":
		deleteSubCommand.Parse(os.Args[2:])
		deleteCommand(taskService, taskId)
	default:
		log.Printf("Comando `%s` no soportado\n", os.Args[1])
	}

}
