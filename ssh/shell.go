package ssh

import (
	"fmt"
)

/**
    Shell creates a shell connection to the host and allows the user to run commands on the server
    @Todo Fix issue using arrow keys in a shell connection
 */
func (client *Client) Shell() error {
	session, err := client.StartSession(true, true)
	if err != nil {
		return fmt.Errorf("Unable to contact server[%s]: %s", client.Config.Host, err)
	}
	defer session.Close()

	// Start remote shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %s", err)
	}

	// wait for the ssh connection the be created and for the shell to stop before continuing
	// Wait for the SCP connection to close, meaning it has consumed all
	// our data and has completed. Or has errored.

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}
