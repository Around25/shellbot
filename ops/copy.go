package ops

import (
	"fmt"
	"github.com/Around25/shellbot/ssh"
	"github.com/mitchellh/go-homedir"
)

func Copy(fromWithHost, toWithHost string, appConfig *Config) error {
	fromHost, fromPath := SplitIdentifierFromPath(fromWithHost)
	toHost, toPath := SplitIdentifierFromPath(toWithHost)

	if fromHost != "" && toHost != "" {
		return fmt.Errorf("Unable to copy between servers yet... Sorry")
	}

	var name string
	if fromHost != "" {
		name = fromHost
	} else {
		name = toHost
	}
	// load connection data from the configuration file
	server, err := appConfig.GetServer(name)
	if err != nil {
		return err
	}
	uri := server["uri"]
	key, _ := homedir.Expand(server["key"])

	// setup new connection to server
	sshConfig := ssh.NewConfig(uri, key, true, true, true)
	client := ssh.New(sshConfig)
	defer client.Disconnect()

	err = client.TryConnection()
	if err != nil {
		return err
	}

	// start copy data transfer
	if toHost != "" {
		err = client.Copy(fromPath, toPath)
	} else {
		err = client.Download(fromPath, toPath)
	}

	if err != nil {
		return fmt.Errorf("Unable to copy from %s to %s: %s\n", fromWithHost, toWithHost, err)
	}
	return nil
}
