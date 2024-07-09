package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	configCmd.AddCommand(configsSetCmd)
}

var configsSetCmd = &cobra.Command{
	Use:   "set config_name config_value",
	Short: "Set configurations for the cli tool",
	Long:  `Set configurations for the cli tool.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		initConfigs()

		configName := args[0]
		configValue := args[1]
		viper.Set(configName, configValue)

		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
	},
}
