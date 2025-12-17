package esxi

import (
	"fmt"
	"log"
)

type Config struct {
	esxiHostName       string
	esxiHostSSHport    string
	esxiHostSSLport    string
	esxiUserName       string
	esxiPassword       string
	esxiPrivateKeyPath string
}

func (c *Config) validateEsxiCreds() error {
	esxiConnInfo := getConnectionInfo(c)
	log.Printf("[validateEsxiCreds]\n")

	var remote_cmd string
	var err error

	remote_cmd = "vmware --version"
	_, err = runRemoteSshCommand(esxiConnInfo, remote_cmd, "Connectivity test, get vmware version")
	if err != nil {
		return fmt.Errorf("failed to connect to esxi host: %s", err)
	}

	runRemoteSshCommand(esxiConnInfo, "mkdir -p ~", "Create home directory if missing")

	return nil
}
