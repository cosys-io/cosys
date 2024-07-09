package cms

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate <command>",
	Short: "Generate code",
	Long:  "Generate code.",
	Run:   func(cmd *cobra.Command, args []string) {},
}
