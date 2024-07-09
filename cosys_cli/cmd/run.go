package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run command_name [arguments]",
	Short: "Run a command from the project binary",
	Long:  "Run a command from the project binary.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initConfigs()

		if !viper.InConfig("bin_path") {
			log.Fatal("configuration not found: bin_path")
		}
		binPath := viper.GetString("bin_path")
		exists, err := pathExists(binPath)
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			log.Fatal("binary path does not exist")
		}

		command := strings.Join(append([]string{binPath}, args...), " ")
		if err := RunCommand(command); err != nil {
			log.Fatal(err)
		}
	},
}
