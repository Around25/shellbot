package cmd

import (
	"github.com/Around25/shellbot/ops"
	"github.com/Around25/shellbot/logger"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy a file or directory between the local environment and a specified server",
	Long: `Copy a file or directory between the local environment and a specified server`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ops.Copy(args[0], args[1], ops.NewConfig(AppConfig))
		if err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(copyCmd)
}
