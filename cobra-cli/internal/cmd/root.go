package cmd

import (
	"cobra-cli/internal/task"
	"context"
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

func Execute(ctx context.Context, s *task.Service) {
	service = s
	if err := rootCommand.ExecuteContext(ctx); err != nil {
		fmt.Println("Error al ejecutar el cli")
		os.Exit(1)
	}
}
