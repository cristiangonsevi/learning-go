package main

import (
	"log"

	"github.com/fatih/color"
)

var availableProfiles []string

func main() {
	var profile string = "qa"
	logger.Debug("Hola debug")
	config, err := loadConfigFile("config.toml")
	if err != nil {
		log.Fatal("Error loading configuration ", err)
	}

	err = validateProfile(profile, config)
	if err != nil {
		log.Fatal(err)
	}

	p := loadProfile(profile, config)
	color.Blue("✅ Profile %s loaded", profile)
	client, err := clientSSH(p)
	if err != nil {
		log.Fatal("Error creating client ", err)
	}
	defer client.Close()
	go keepAlive(client)
	for _, s := range p.Services {
		if !s.Enabled {
			color.Yellow("⚙️  %s disabled in config file\n", s.Name)
			continue
		}
		ready := make(chan error, 1)
		go startTunnel(client, s, ready)

		if err := <-ready; err != nil {
			color.Red("✗ %s failed: %v\n", s.Name, err)
			continue
		}
		color.Green("⚙️  %s started on port %d -> %d\n", s.Name, s.LocalPort, s.RemotePort)
	}
	select {}
}
