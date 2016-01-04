package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the status of your servers and their current state",
	Long: `Check the status of your servers and their current state`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Check command not yet implemented!!!")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
