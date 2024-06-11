package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os/exec"
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
	cmd := exec.Command("./bin/cosys")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
