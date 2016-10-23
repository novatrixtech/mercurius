package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of Mercurius",
	Long: `Version of Mercurius`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mercurius version", VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
