package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cristiangonsevi/learning-go/portbridge/internal/logger"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

const sampleConfiguration = `# PortBridge — configuración de ejemplo.
# Edita este archivo y luego corre: portbridge <perfil>

[profiles.sample]
user       = "your_user"
host       = "your_host:22"
knownhost  = "~/.ssh/known_hosts"
keyfile    = "~/.ssh/id_ed25519"
passphrase = ""
password   = ""

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

var ConfigurationPath string

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

func InitDefaultConfig(forceRecreate bool) {
	defaultLocation := getDefaultConfigFilePath()
	_, err := os.ReadFile(defaultLocation + "/config.toml")
	if err == nil && !forceRecreate {
		color.Yellow("Configuration file already exist on %s\n", defaultLocation)
		return
	}

	err = os.MkdirAll(defaultLocation, 0o755)
	if err != nil {
		logger.Error("Error creating directory %v", err)
	}

	f, err := os.Create(defaultLocation + "/config.toml")
	if err != nil {
		log.Fatal("Error creating default configuration file ", err)
	}

	_, err = f.WriteString(sampleConfiguration)
	if err != nil {
		log.Fatal("Error writing default configuration to file ", err)
	}
	color.Green("Default configuration file created in %v\n%v", defaultLocation, sampleConfiguration)
}

func ValidateProfile(profile string, config *Config) error {
	_, ok := config.Profiles[profile]

	if !ok {
		return fmt.Errorf(color.RedString("Profile %s not exist", profile))
	}
	return nil
}

func LoadProfile(profile string, config *Config) Profile {
	p, ok := config.Profiles[profile]
	if !ok {
		fmt.Errorf("Profile not exit")
		return Profile{}
	}
	return p
}

func LoadConfigFile() (*Config, error) {
	if ConfigurationPath == "" {
		ConfigurationPath = getDefaultConfigFilePath() + "/config.toml"
	}
	var config *Config
	_, err := toml.DecodeFile(ConfigurationPath, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SetConfigurationPath(path string) {
	ConfigurationPath = path
}

func GetAvailableProfiles(config *Config) []string {
	var profilesName []string
	for k := range config.Profiles {
		profilesName = append(profilesName, k)
	}
	return profilesName
}

func getDefaultConfigFilePath() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error loading default config directory", err)
	}
	return filepath.Join(userConfigDir, "portbridge")
}
