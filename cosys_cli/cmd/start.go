package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

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

func startServer() error {
	initConfigs()

	binPath := viper.GetString("bin_path")
	exists, err := pathExists(binPath)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Fatal("binary path does not exist")
	}

	if err = RunCommand(binPath + " serve"); err != nil {
		return err
	}

	return nil
}
