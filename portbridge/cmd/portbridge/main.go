package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cristiangonsevi/learning-go/portbridge/internal/cli"
)

var availableProfiles []string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cli.Execute(ctx)
	// var profile string = "qa"
	// logger.Debug("Hola debug")
	// cfg, err := config.LoadConfigFile("config.toml")
	// if err != nil {
	// 	logger.Error("Error loading configuration %v", err)
	// 	return
	// }
	//
	// err = config.ValidateProfile(profile, cfg)
	// if err != nil {
	// 	logger.Error("Error validating profile %v", err)
	// 	return
	// }
	//
	// p := config.LoadProfile(profile, cfg)
	// color.Blue("✅ Profile %s loaded", profile)
	// client, err := ssh.ClientSSH(p)
	// if err != nil {
	// 	log.Fatal("Error creating client ", err)
	// }
	// defer client.Close()
	// go ssh.KeepAlive(client)
	// for _, s := range p.Services {
	// 	if !s.Enabled {
	// 		color.Yellow("⚙️  %s disabled in config file\n", s.Name)
	// 		continue
	// 	}
	// 	ready := make(chan error, 1)
	// 	go tunnel.StartTunnel(client, s, ready)
	//
	// 	if err := <-ready; err != nil {
	// 		color.Red("✗ %s failed: %v\n", s.Name, err)
	// 		continue
	// 	}
	// 	color.Green("⚙️  %s started on port %d -> %d\n", s.Name, s.LocalPort, s.RemotePort)
	// }
	// select {}
}
