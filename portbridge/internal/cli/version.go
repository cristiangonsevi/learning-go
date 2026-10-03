package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display current version of portbridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Portbridge v%v\n", version)
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
