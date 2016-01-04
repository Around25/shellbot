package cmd

import (
	"github.com/Around25/shellbot/ops"
	"github.com/Around25/shellbot/logger"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Execute all the tasks associated with each group of servers",
	Long: `Execute all the tasks associated with each group of servers`,
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) == 0 {
			logger.Fatal("No environment specified")
		}
		name = args[0]

		err := ops.ProvisionEnvironment(name, ops.NewConfig(AppConfig))
		if err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
