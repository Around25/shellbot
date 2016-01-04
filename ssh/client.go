package ssh

import (
	"fmt"
	"io"
	"os"
	"golang.org/x/crypto/ssh"
)

/**
	Client holds an active connection to a SSH server
 */
type Client struct {
	Config *Config
	conn   *ssh.Client
}

/**
	New creates a new client connection to the server specified by the configuration
 */
func New(config *Config) (client *Client) {
	client = &Client{
		Config: config,
	}
	return client
}

/**
	Connect the client to the configured server
 */
func (client *Client) Connect() (error) {
	config, err := client.Config.GetAuthConfig()
	if err != nil {
		return err
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Config.Host, client.Config.Port), config)
	if err != nil {
		return fmt.Errorf("Failed to dial: %s", err)
	}
	client.conn = connection
	return nil
}

func (client *Client) Disconnect() {
	if client.conn != nil {
		client.conn.Close()
	}
}

/**
    Connect to the server base on a client auth config
 */
func (client *Client) StartSession(bindIOStreams bool, createPty bool) (*ssh.Session, error) {
	// make sure the connection is available
	if (client.conn == nil) {
		err := client.Connect()
		if err != nil {
			return nil, err
		}
	}
	// start a new SSH session
	session, err := client.conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	// bind IO streams
	if (bindIOStreams) {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin
	}

	// create PTY
	if (createPty) {
		modes := ssh.TerminalModes{
			ssh.ECHO:          0, // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			session.Close()
			return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
		}
	}

	return session, nil
}

/**
	Try the connection to the server
 */
func (client *Client) TryConnection() error {
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("Unable to contact server[%s]: %s", client.Config.Host, err)
	}
	return nil
}

/**
    InitIOPipes binds session IO streams to the output
    @Deprecated
 */
func (client *Client) InitIOPipes(session *ssh.Session) error {
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)
	return nil
}
