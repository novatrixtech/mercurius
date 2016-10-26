package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Check Mercurius Commands",
	Long:  `Check Mercurius Commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// print help usage string
		fmt.Println(RootCmd.UsageString())
	},
}

func init() {
	RootCmd.AddCommand(helpCmd)
}
