package internal

import "github.com/spf13/cobra"

// RootCmd is the root command of the cli tool commands for the cms.
var RootCmd = &cobra.Command{
	Use:   "cms <command>",
	Short: "Manage the cms.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
