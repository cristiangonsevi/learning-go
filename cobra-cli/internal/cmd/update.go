package cmd

import (
	"cobra-cli/internal/cmd/printer"
	"cobra-cli/internal/model"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var updateCommand = cobra.Command{
	Use:   "update",
	Short: "Comando para actualizar nombre o estado de tarea",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		id, _ := cmd.Flags().GetInt("id")

		name, _ := cmd.Flags().GetString("name")

		status, _ := cmd.Flags().GetString("status")

		if name == "" && status == "" {
			fmt.Println("Debes proporcionar los valores de al menos nombre o estado")
			cmd.Help()
			return
		}

		var taskItem = model.TaskModel{
			Name:   name,
			Status: status,
		}

		err := service.UpdateTask(ctx, id, taskItem)

		if err != nil {
			log.Fatal("Error actualizando tarea ", err)
		}

		tasks, _ := service.List(ctx)
		printer.New().PrintTask(tasks)
	},
}

func init() {

	updateCommand.Flags().Int("id", 0, "Id de la tarea a actualizar")

	updateCommand.MarkFlagRequired("id")

	updateCommand.Flags().StringP("name", "n", "", "Nombre con el que se actualizara la tarea")

	updateCommand.Flags().StringP("status", "s", "", "Estatus a actualizar de la tarea")

	rootCommand.AddCommand(&updateCommand)
}
