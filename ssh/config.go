package ssh

import (
	"github.com/Around25/shellbot/logger"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config holds the configuration properties for one server connection
type Config struct {
	// connection properties
	Host string
	Port int

	// authentication properties
	User     string
	Password string
	AuthFile string
	SSHAgent bool

	// session properties
	CreatePty     bool
	BindIOStreams bool

	// active connection properties
	Config *ssh.ClientConfig
}

// New Config creates a new Config object based on the given URI and other data
func NewConfig(uri string, authFile string, sshAgent bool, createPty bool, bindIoStreams bool) *Config {
	// prefix the uri with ssh:// if invalid
	if !strings.HasPrefix(uri, "ssh://") {
		uri = "ssh://" + uri
	}
	parsed, err := url.Parse(uri)
	if err != nil {
		logger.Fatal("Invalid uri provided: " + uri)
	}

	var (
		user string
		pass string
		port string
		host string
	)

	// load the user and password if provided
	if parsed.User != nil {
		user = parsed.User.Username()
		pass, _ = parsed.User.Password()
	}
	// load host and port if available
	host, port, err = net.SplitHostPort(parsed.Host)
	if err != nil {
		host = parsed.Host
	}

	// add default values if empty
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "22"
	}

	iPort, _ := strconv.Atoi(port)

	// return the filled config object
	return &Config{
		Host:          host,
		User:          user,
		Password:      pass,
		Port:          iPort,
		AuthFile:      authFile,
		SSHAgent:      sshAgent,
		CreatePty:     createPty,
		BindIOStreams: bindIoStreams,
	}
}

/**
  GetAuthConfig loads an authentication configuration based on the user and auth method provided
*/
func (config *Config) GetAuthConfig() (*ssh.ClientConfig, error) {
	var auth ssh.AuthMethod = nil
	if len(config.Password) != 0 {
		auth = config.AuthViaPassword(config.Password)
	} else if len(config.AuthFile) != 0 {
		auth = config.AuthViaKey(config.AuthFile)
	} else if config.SSHAgent {
		auth = config.AuthViaSSHAgent()
	}

	config.Config = &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{auth},
	}
	return config.Config, nil
}

/**
  AuthViaPassword returns an authentication method using a password as a credential
*/
func (config *Config) AuthViaPassword(pass string) ssh.AuthMethod {
	return ssh.Password(pass)
}

/**
  AuthViaKey returns an authentication method using a private credential file
*/
func (config *Config) AuthViaKey(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

/**
  AuthViaSSHAgent returns an authentication method using an SSH Agent with loaded keys
*/
func (config *Config) AuthViaSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
