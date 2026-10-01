package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

type Service struct {
	Name       string `toml:"name"`
	RemotePort int    `toml:"remoteport"`
	LocalPort  int    `toml:"localport"`
	Enabled    bool   `toml:"enabled"`
}
type Profile struct {
	Name       string    `toml:"name"`
	User       string    `toml:"user"`
	Host       string    `toml:"host"`
	KnownHost  string    `toml:"knownhost"`
	KeyFile    string    `toml:"keyfile"`
	PassPhrase string    `toml:"passphrase"`
	Password   string    `toml:"password"`
	Services   []Service `toml:"services"`
}

type Config struct {
	Profiles map[string]Profile `toml:"profiles"`
}

func initDefaultConfig() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error loading default config directory", err)
	}
	defaultLocation := userConfigDir + "/portbridge/config.toml"
	_, err = os.ReadFile(defaultLocation)
	if err == nil {
		color.Yellow("Loading default configuration file %s\n", defaultLocation)
		return
	}

	f, err := os.Create(defaultLocation)
	if err != nil {
		log.Fatal("Error creating default configuration file ", err)
	}

	sampleConfiguration := `# PortBridge — configuración de ejemplo.
# Edita este archivo y luego corre: portbridge <perfil>

[profiles.sample]
user       = "your_user"
host       = "your_host:22"
knownhost  = "~/.ssh/known_hosts"
keyfile    = "~/.ssh/id_ed25519"
passphrase = ""

[[profiles.sample.services]]
name       = "redis"
remotehost = "localhost"
remoteport = 6379
localport  = 6379
enabled    = true

[[profiles.sample.services]]
name       = "postgres"
remotehost = "localhost"
remoteport = 5432
localport  = 5432
enabled    = true
`
	_, err = f.WriteString(sampleConfiguration)
	if err != nil {
		log.Fatal("Error writing default configuration to file ", err)
	}
	color.Green("Default configuration file created in %v\n%v", defaultLocation, sampleConfiguration)
}

func validateProfile(profile string, config *Config) error {
	availableProfiles := getAvailableProfiles(config)
	for _, p := range availableProfiles {
		if p == profile {
			return nil
		}
	}
	return fmt.Errorf(color.RedString("Profile %s not exist, available profiles: \n%v", profile, strings.Join(availableProfiles, "\n")))
}

func loadProfile(profile string, config *Config) Profile {
	p, ok := config.Profiles[profile]
	if !ok {
		fmt.Errorf("Profile not exit")
		return Profile{}
	}
	return p
}

func loadConfigFile(path string) (*Config, error) {
	var config *Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getAvailableProfiles(config *Config) []string {
	var profilesName []string
	for k := range config.Profiles {
		profilesName = append(profilesName, k)
	}
	return profilesName
}
