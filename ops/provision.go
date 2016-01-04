package ops
import (
	"fmt"
	"github.com/Around25/shellbot/ssh"
	"github.com/mitchellh/go-homedir"
)

/**
	Provision an environment based on the loaded configuration file
 */
func ProvisionEnvironment(env string, config *Config) error {
	groups := config.GetGroupsForEnv(env)
	variables := config.GetVariablesForEnv(env)

	for _, group := range groups {
		err := ProvisionGroup(group, variables, config)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
	Provision a specific group from the config using the given variables
 */
func ProvisionGroup(group string, variables map[string] string, config *Config) error {
	servers := config.GetServersForGroup(group)
	tasks := config.GetTasksForGroup(group)
	checks := config.GetChecksForGroup(group)

	for _, server := range servers {
		err := ProvisionServer(server, tasks, checks, variables, config)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
	Provision a single server with the list of tasks and variables
 */
func ProvisionServer(name string, tasks []map[string] string, checks []map[string] string, variables map[string] string, config *Config) error {
	// load connection data from the configuration file
	server, err := config.GetServer(name)
	if err != nil {
		return err
	}
	uri := server["uri"];
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

	// execute tasks on the current server
	err = ExecuteTasksOnServer(client, tasks, variables, config)
	if err != nil {
		return err
	}

	// @todo Execute checks on server
	return nil
}

/**
	Execute a list of tasks on the connected server based on the provided config and with the list of variables as the current context
 */
func ExecuteTasksOnServer(client *ssh.Client, tasks []map[string] string, variables map[string] string, config *Config) error {
	for _, task := range tasks {
		for taskType, taskValue := range task {
			taskValue = ExpandVariables(taskValue, variables)
			output, err := ExecuteTaskOnServer(client, taskType, taskValue, config)
			if err != nil {
				return err
			}
			fmt.Print(output)
		}
	}
	return nil
}

/**
	Execute the current task on the server
 */
func ExecuteTaskOnServer(client *ssh.Client, taskType string, taskValue string, config *Config) (string, error) {
	switch taskType {
	case "run":
		return client.Execute(taskValue)
	case "task":
		return ExecuteTaskGroupOnServer(client, taskValue, config)
	case "copy":
		from, to := splitPaths(taskValue)
		return "", client.Copy(from, to)
	case "download":
		from, to := splitPaths(taskValue)
		return "", client.Download(from, to)
	}
	return "", fmt.Errorf("Unknown task type: %s", taskType)
}

/**
	Execute a task group on the server
 */
func ExecuteTaskGroupOnServer(client *ssh.Client, group string, config *Config) (string, error) {
	tasks := config.GetTasksForTaskGroup(group)
	err := ExecuteTasksOnServer(client, tasks, nil, config)
	if err != nil {
		return "", err
	}
	return "", nil
}