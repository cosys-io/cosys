package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the cli tool.
var rootCmd = &cobra.Command{
	Use:   "cosys <command> [arguments]",
	Short: "Manage a cosys project.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	}}

// Execute runs the cli tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
