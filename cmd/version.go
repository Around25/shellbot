package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of the tool",
	Long:  `Display the current version of the tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.1.0")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
