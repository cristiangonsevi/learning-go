package main

import (
	"fmt"
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
	Profiles []Profile `toml:"profiles"`
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
	for _, p := range config.Profiles {
		if p.Name == profile {
			return p
		}
	}
	return Profile{}
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
	for _, p := range config.Profiles {
		profilesName = append(profilesName, p.Name)
	}
	return profilesName
}
