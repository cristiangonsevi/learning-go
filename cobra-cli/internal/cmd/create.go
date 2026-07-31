package cmd

import (
	"cobra-cli/internal/cmd/printer"
	"cobra-cli/internal/model"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var createCommand = cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")

		newTask := model.TaskModel{
			Name:      name,
			Status:    status,
			CreatedAt: time.Now(),
		}

		err := service.CreateTask(newTask)

		if err != nil {
			log.Fatal("Error creating new task ", err)
		}

		tasks, _ := service.List()

		printer.New().PrintTask(tasks)
	},
}

func init() {
	createCommand.Flags().StringP("name", "n", "", "Define nombre de la tarea a crear")
	createCommand.MarkFlagRequired("name")
	createCommand.Flags().StringP("status", "s", "pending", "Define status de la tarea")
	rootCommand.AddCommand(&createCommand)
}
