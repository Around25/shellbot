package ops

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// config contains the contents of the loaded configuration file and provides methods for easily retrieving it's data
type Config struct {
	config *viper.Viper
}

/**
Create a new Config object based on a viper config
*/
func NewConfig(config *viper.Viper) *Config {
	return &Config{
		config: config,
	}
}

/**
Retrieve the connection details for a specific server
*/
func (config *Config) GetServer(name string) (map[string]string, error) {
	server := config.config.GetStringMapString("servers." + name)
	if server == nil {
		return nil, fmt.Errorf("Server is not defined")
	}
	return server, nil
}

/**
Retrieve the list of groups for the specified environment
*/
func (config *Config) GetGroupsForEnv(env string) []string {
	return config.config.GetStringSlice("environments." + env + ".groups")
}

/**
Retrieve the list of variables for the specified environment from the config file and from the OS
*/
func (config *Config) GetVariablesForEnv(env string) map[string]string {
	variables := config.config.GetStringMapString("environments." + env + ".variables")
	for key, _ := range variables {
		if os.Getenv(key) != "" {
			variables[key] = os.Getenv(key)
		}
	}
	return variables
}

/**
Retrieve the list of servers contained in a group
*/
func (config *Config) GetServersForGroup(group string) []string {
	return config.config.GetStringSlice("groups." + group + ".servers")
}

/**
Retrieve the list of tasks that should be executed for a group
*/
func (config *Config) GetTasksForGroup(group string) []map[string]string {
	// retrieve the list from the config file
	tasks := config.config.Get("groups." + group + ".tasks")

	// convert from interface{} to []map[string] string and return
	return convertFromInterface(tasks)
}

/**
Retrieve the list of tasks included in a task group
*/
func (config *Config) GetTasksForTaskGroup(taskGroup string) []map[string]string {
	// retrieve the list from the config file
	tasks := config.config.Get("tasks." + taskGroup)

	// convert from interface{} to []map[string] string and return
	return convertFromInterface(tasks)
}

/**
Retrieve the list of checks for a specific group
*/
func (config *Config) GetChecksForGroup(group string) []map[string]string {
	checks := config.config.Get("groups." + group + ".checks")

	// convert from interface{} to []map[string] string and return
	return convertFromInterface(checks)
}

/**
Convert an interface from Viper to an []map[string] string type
*/
func convertFromInterface(rawData interface{}) []map[string]string {
	data := rawData.([]interface{})
	result := make([]map[string]string, len(data))
	for k, v := range data {
		val := v.(map[interface{}]interface{})
		item := make(map[string]string, len(val))
		for t, d := range val {
			item[t.(string)] = d.(string)
		}
		result[k] = item
	}
	return result
}
