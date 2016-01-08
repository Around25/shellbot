package ops

import (
	"fmt"

	"github.com/Around25/shellbot/ssh"
	"github.com/mitchellh/go-homedir"
)

func OpenTerminalToServer(name string, appConfig *Config) error {
	var (
		sshConfig *ssh.Config
		uri       string
		key       string
	)

	server, _ := appConfig.GetServer(name)
	if server != nil {
		uri = server["uri"]
		key, _ = homedir.Expand(server["key"])
	}

	if uri == "" {
		sshConfig = ssh.NewConfig(name, "", false, true, true)
	} else {
		// load connection data from the configuration file
		sshConfig = ssh.NewConfig(uri, key, true, true, true)

	}

	client := ssh.New(sshConfig)
	defer client.Disconnect()

	// test the connection
	err := client.TryConnection()
	if err != nil {
		return fmt.Errorf("Unable to connect to server using uri '%s' key '%s' got: %s", uri, key, err)
	}

	// start the terminal shell
	err = client.Shell()
	if err != nil {
		return err
	}
	return nil
}
