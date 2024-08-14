package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

// initConfigs loads the configurations from the .cli_configs file.
func initConfigs() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName(".cli_configs")
	viper.AddConfigPath(dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// configCmd is the command for managing configurations for the cli tool.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations for the cli tool",
	Long:  "Manage configurations for the cli tool.",
	Run:   func(cmd *cobra.Command, args []string) {},
}
