package ops

import (
	"fmt"
	"github.com/Around25/shellbot/ssh"
	"github.com/mitchellh/go-homedir"
)

func OpenTerminalToServer(name string, appConfig *Config) error {
	// load connection data from the configuration file
	server, err := appConfig.GetServer(name)
	if err != nil {
		return err
	}
	uri := server["uri"]
	key, _ := homedir.Expand(server["key"])

	if uri == "" {
		return fmt.Errorf("Missing connection string. Define your server in the configuration file first.")
	}

	// create a new client
	sshConfig := ssh.NewConfig(uri, key, true, true, true)
	client := ssh.New(sshConfig)
	defer client.Disconnect()

	// test the connection
	err = client.TryConnection()
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
