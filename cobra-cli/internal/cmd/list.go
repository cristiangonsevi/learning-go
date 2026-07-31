package cmd

import (
	"cobra-cli/internal/cmd/printer"
	"log"

	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Lista todos",
	Long:  "Muestra el listado de comandos",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := service.List()
		if err != nil {
			log.Fatal("Error retrieving data ", err)
		}

		printer.New().PrintTask(tasks)

	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
