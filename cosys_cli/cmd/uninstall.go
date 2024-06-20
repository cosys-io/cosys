package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [module_names]",
	Short: "Uninstall modules",
	Long:  "Uninstall modules.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, moduleName := range args {
			if err := os.RemoveAll(filepath.Join("modules", moduleName)); err != nil {
				log.Fatal(err)
			}
		}
	},
}
