package ssh

import (
    "fmt"
)

/**
    Execute a command on the server
 */
func (client *Client) Execute(command string) (string, error) {
    session, err := client.StartSession(false, true)
    if err != nil {
        return "", fmt.Errorf("Unable to contact server[%s]: %s", client.Config.Host, err)
    }
    defer session.Close()

//	session.Run(command) // without capturing the output, just the error; the output can be banded to a io.Writer
    output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), err
	}

    return string(output), nil
}