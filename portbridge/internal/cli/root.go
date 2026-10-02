package cli

import (
	"context"

	"portbridge/internal/config"
	"portbridge/internal/logger"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "pl",
	Short: "Create SSH tunnels for local development",
	Long:  "Create secure SSH tunnels to access remote services from your local environment without exposing them to the public internet.",
	Run: func(cmd *cobra.Command, args []string) {
		customConfigPath, err := cmd.Flags().GetString("config")
		if err != nil {
			logger.Error("Error getting config flag %v", err)
			return
		}

		if customConfigPath != "" {
			logger.Info("Custom path provided %s", customConfigPath)
		}
		cmd.Help()
	},
}

func Execute(ctx context.Context) {
	rootCommand.PersistentFlags().StringVarP(&config.ConfigurationPath, "config", "c", "", "Custom configuration file path")
	if err := rootCommand.ExecuteContext(ctx); err != nil {
		logger.Error("Error executing this program %v", err)
		return
	}
}
