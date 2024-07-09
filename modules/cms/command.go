package cms

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "cms <command>",
	Short: "Manage the cms.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
