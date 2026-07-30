package cmd

import (
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Lista todos",
	Long:  "Muestra el listado de comandos",
	Run: func(cmd *cobra.Command, args []string) {
		service.List()
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
