package cli

import (
	"github.com/cristiangonsevi/learning-go/portbridge/cmd/pb/internal/config"
	"github.com/cristiangonsevi/learning-go/portbridge/cmd/pb/internal/logger"

	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Create default configuration file",
	Long:  "Create a default configuration file with a sample",
	Run: func(cmd *cobra.Command, args []string) {
		forceRecreate, err := cmd.Flags().GetBool("force")
		if err != nil {
			logger.Error("Error getting flag %v", err)
		}
		config.InitDefaultConfig(forceRecreate)
	},
}

func init() {
	initCommand.Flags().BoolP("force", "f", false, "Force recreation of sample configuration")
	rootCommand.AddCommand(initCommand)
}
