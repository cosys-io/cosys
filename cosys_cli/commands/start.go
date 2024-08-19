package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

// startCmd is the command for starting the server.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Deploy the Golang server and the Content Management UI server",
	Long:  `Deploy the Golang server and the Content Management UI server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := startServer(); err != nil {
			log.Fatal(err)
		}
	},
}

// startServer starts the server.
func startServer() error {
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

	if err = runCommand(binPath + " serve"); err != nil {
		return err
	}

	return nil
}
