package cmd

import (
	"cobra-cli/internal/task"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var service *task.Service

var rootCommand = &cobra.Command{
	Use:   "cobra-cli",
	Short: "Short description",
	Long:  "Long description",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(s *task.Service) {
	service = s
	if err := rootCommand.Execute(); err != nil {
		fmt.Println("Error al ejecutar el cli")
		os.Exit(1)
	}
}
