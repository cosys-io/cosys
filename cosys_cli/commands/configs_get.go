package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	configCmd.AddCommand(configsGetCmd)
}

var configsGetCmd = &cobra.Command{
	Use:   "get cfg_name",
	Short: "Get configurations for the cli tool",
	Long:  `Get configurations for the cli tool.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initConfigs()

		cfgName := args[0]
		if !viper.InConfig(cfgName) {
			log.Fatalf("cfg not found: %s", cfgName)
		}

		cfgValue := viper.Get(cfgName)
		log.Print(cfgValue)
	},
}
