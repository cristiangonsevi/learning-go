package cli

import (
	"log"
	"strings"

	"portbridge/internal/config"
	"portbridge/internal/logger"
	"portbridge/internal/ssh"
	"portbridge/internal/tunnel"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var upCommand = &cobra.Command{
	Use:   "up <profile>",
	Short: "Activate an SSH tunnel profile",
	Long:  `Activate an SSH tunnel profile and establish connections to its configured services.`,
	Example: `
	pb up qa
	pb up qa --services redis,postgres
	pb up qa -s redis,postgres`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		pickedServices, err := cmd.Flags().GetString("services")
		if err != nil {
			logger.Error("Error getting services flag %v", err)
			return
		}

		profile := args[0]
		cfg, err := config.LoadConfigFile()
		if err != nil {
			logger.Error("Error loading configuration file %v", err)
			return
		}
		err = config.ValidateProfile(profile, cfg)
		if err != nil {
			logger.Error("%v", err)
			printProfiles()
			return
		}

		p := config.LoadProfile(profile, cfg)
		color.Blue("✅ Profile %s loaded", profile)
		client, err := ssh.ClientSSH(p)
		if err != nil {
			log.Fatal("Error creating client ", err)
		}
		defer client.Close()
		go ssh.KeepAlive(client)
		for _, s := range p.Services {
			matchService := strings.Contains(strings.ToLower(pickedServices), strings.ToLower(s.Name))
			if pickedServices != "" && !matchService {
				continue
			}
			if !s.Enabled && !matchService {
				color.Yellow("⚙️  %s disabled in config file\n", s.Name)
				continue
			}
			ready := make(chan error, 1)
			go tunnel.StartTunnel(client, s, ready)

			if err := <-ready; err != nil {
				color.Red("✗ %s failed: %v\n", s.Name, err)
				continue
			}
			color.Green("⚙️  %s started on port %d -> %d\n", s.Name, s.LocalPort, s.RemotePort)
		}
		<-cmd.Context().Done()
	},
}

func init() {
	upCommand.Flags().StringP("services", "s", "", "Comma-separated service names, e.g. redis,postgres")
	rootCommand.AddCommand(upCommand)
}
