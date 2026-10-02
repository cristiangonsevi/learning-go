package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"portbridge/internal/config"
	"portbridge/internal/logger"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List commando to display available configured profiles",
	Run: func(cmd *cobra.Command, args []string) {
		printProfiles()
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}

func printProfiles() {
	cfg, err := config.LoadConfigFile()
	if err != nil {
		logger.Error("Error loading config file %v", err)
		return
	}
	profiles := config.GetAvailableProfiles(cfg)

	fmt.Println("=================================")
	fmt.Println("Available profiles")
	fmt.Println("=================================")
	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{"no", "profile", "usage", "services"})
	for i, v := range profiles {
		var servicesList []string
		profile := config.LoadProfile(v, cfg)
		for _, val := range profile.Services {
			serviceName := val.Name
			if !val.Enabled {
				serviceName += " (disabled)"
			}
			servicesList = append(servicesList, serviceName)
		}
		table.Append([]string{strconv.Itoa(i + 1), v, fmt.Sprintf("pl up %s", v), strings.Join(servicesList, ", ")})
	}
	table.Render()
}
