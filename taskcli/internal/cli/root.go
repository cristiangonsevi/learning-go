package cli

import (
	"flag"
	"fmt"
	"os"
)

func NewRootCommand() {
	createCmd := flag.NewFlagSet("list", flag.ExitOnError)
	taskName := createCmd.String("name", "", "Nombre de la tarea")

	if len(os.Args) < 2 {
		fmt.Printf("Mostrando ayudame de cli\n")
		return
	}

	switch os.Args[1] {
	case "list":
		newListCommand()
	case "create":
		createCmd.Parse(os.Args[2:])
		createCommand(taskName)
	}

}
