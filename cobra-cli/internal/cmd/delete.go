package cmd

import (
	"cobra-cli/internal/cmd/printer"
	"log"

	"github.com/spf13/cobra"
)

var deleteCommand = cobra.Command{
	Use:   "delete",
	Short: "Comando para borrar una tarea",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		id, _ := cmd.Flags().GetInt("id")

		err := service.DeleteTask(ctx, id)

		if err != nil {
			log.Fatal("Error al borrar tarea ", err)
		}

		tasks, _ := service.List(ctx)

		printer.New().PrintTask(tasks)
	},
}

func init() {
	deleteCommand.Flags().Int("id", 0, "Id de la tarea a borrar")
	deleteCommand.MarkFlagRequired("id")

	rootCommand.AddCommand(&deleteCommand)
}
