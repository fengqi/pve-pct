package client

import (
	"fmt"
	"log"

	"github.com/fengqi/pve-pct/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	ssh    *ssh.Client
	sftp   *sftp.Client
	Config *config.Config
	Cmd    string
}

func Init(c *config.Config) *SSHClient {
	signer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(c.PrivateKey), []byte(c.Password))
	if err != nil {
		log.Fatal("failed to parse private key")
	}

	scc := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: VerifyHost,
	}

	s, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Addr, c.Port), scc)
	if err != nil {
		log.Fatalf("failed to dial pve %s:%d, err %v", c.Addr, c.Port, err)
	}

	sf, err := sftp.NewClient(s)
	if err != nil {
		log.Fatalf("failed to dial pve %s:%d, err %v", c.Addr, c.Port, err)
	}

	return &SSHClient{
		Config: c,
		ssh:    s,
		sftp:   sf,
	}
}
