package ssh

import (
	"fmt"
	"os"
	"time"

	"github.com/cristiangonsevi/learning-go/portbridge/internal/config"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func KeepAlive(client *ssh.Client) {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
		if err != nil {
			return
		}
	}
}

func ClientSSH(profile config.Profile) (*ssh.Client, error) {
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
