package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var availableProfiles []string

type Service struct {
	Name       string
	RemotePort int
	LocalPort  int
	Enabled    bool
}
type Profile struct {
	Name       string
	User       string
	Host       string
	KnownHost  string
	KeyFile    string
	PassPhrase string
	Password   string
	Services   []Service
}

type Config struct {
	Profiles []Profile
}

func main() {
	var profile string = "qa2"
	config := Config{
		Profiles: []Profile{{
			Name:       "profile",
			User:       "user",
			Host:       "localhost:port",
			KnownHost:  "known_hosts",
			KeyFile:    "datawagon_vps",
			PassPhrase: "passphrase",
			Password:   "",
			Services: []Service{{
				Name:       "Redis",
				RemotePort: 6379,
				LocalPort:  8000,
				Enabled:    false,
			}, {
				Name:       "Postgres",
				RemotePort: 5432,
				LocalPort:  8001,
				Enabled:    true,
			}},
		}},
	}

	err := validateProfile(profile, config)
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
		color.Green("⚙️  %s started on port %d\n", s.Name, s.LocalPort)
	}
	select {}
}

func clientSSH(profile Profile) (*ssh.Client, error) {
	hostKeyCallback, err := knownhosts.New(profile.KnownHost)
	if err != nil {
		return nil, err
	}

	sshConfig := ssh.ClientConfig{
		User:            profile.User,
		HostKeyCallback: hostKeyCallback,
		Auth:            []ssh.AuthMethod{},
	}

	if profile.KeyFile != "" {
		privateKey, err := os.ReadFile(profile.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("error reading key file ", err)
		}

		var signer ssh.Signer

		if profile.PassPhrase != "" {
			signerResponse, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(profile.PassPhrase))
			if err != nil {
				return nil, fmt.Errorf("error parsing key file ", err)
			}

			signer = signerResponse
		} else {
			signerResponse, err := ssh.ParsePrivateKey(privateKey)
			if err != nil {
				return nil, fmt.Errorf("error parsing key file ", err)
			}

			signer = signerResponse
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}

	client, err := ssh.Dial("tcp", profile.Host, &sshConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating ssh client ", err)
	}

	return client, nil
}

func startTunnel(client *ssh.Client, service Service, ready chan<- error) {
	localAddress := fmt.Sprintf("localhost:%d", service.LocalPort)
	listener, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Println("Error listener ", err)
		ready <- err
		return
	}
	defer listener.Close()
	ready <- nil
	for {

		localConn, err := listener.Accept()
		if err != nil {
			log.Println("Error local connections ", err)
			return
		}
		remoteAddres := fmt.Sprintf("localhost:%d", service.RemotePort)
		remoteConn, err := client.Dial("tcp", remoteAddres)
		if err != nil {
			localConn.Close()
			continue
		}

		go func() {
			io.Copy(remoteConn, localConn)
			remoteConn.Close()
		}()
		go func() {
			io.Copy(localConn, remoteConn)
			localConn.Close()
		}()
	}
}

func validateProfile(profile string, config Config) error {
	availableProfiles := getAvailableProfiles(config)
	for _, p := range availableProfiles {
		if p == profile {
			return nil
		}
	}
	return fmt.Errorf(color.RedString("Profile %s not exist, available profiles: \n%v", profile, strings.Join(availableProfiles, "\n")))
}

func getAvailableProfiles(config Config) []string {
	var profilesName []string
	for _, p := range config.Profiles {
		profilesName = append(profilesName, p.Name)
	}
	return profilesName
}

func loadProfile(profile string, config Config) Profile {
	for _, p := range config.Profiles {
		if p.Name == profile {
			return p
		}
	}
	return Profile{}
}
