package ops

import (
	"strings"
	"os"
)

/**
	Split the server name from the path on that server
 */
func SplitIdentifierFromPath(pathWithHost string) (string, string) {
	// the string doesn't contain any : so it's just the path on the host
	if !strings.Contains(pathWithHost, ":") {
		return "", pathWithHost
	}

	parts := strings.SplitN(pathWithHost, ":", 2)
	host := parts[0]
	path := parts[1]
	if host == "localhost" {
		return "", path
	}
	return host, path
}

func splitPaths(data string) (string, string){
	parts := strings.SplitN(data, " ", 2)
	from := parts[0]
	to := parts[1]
	return from, to
}

func ExpandVariables(data string, variables map[string]string) string{
	return os.Expand(data, func (found string) string{
		return variables[found]
	})
}