package cmd

import (
	"github.com/Around25/shellbot/logger"
	"github.com/Around25/shellbot/ops"

	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Connect to a single server through ssh",
	Long:  `Using shell you can connect open a ssh session to a server.`,
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) == 0 {
			logger.Fatal("Specify the name of the server you want to connect to.")
		}

		// get the name of the server as the first argument
		name = args[0]
		err := ops.OpenTerminalToServer(name, ops.NewConfig(AppConfig))
		if err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(shellCmd)
}
