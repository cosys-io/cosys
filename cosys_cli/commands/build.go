package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	mainPath  string // mainPath is bound to the main_path flag.
	indexPath string // indexPath is bound to the index_path flag.
	binPath   string // binPath is bound to the output flag.
)

func init() {
	buildCmd.Flags().StringVarP(&mainPath, "main_path", "M", "", "location of main package")
	buildCmd.Flags().StringVarP(&indexPath, "index_path", "I", "", "location of ui index file")
	buildCmd.Flags().StringVarP(&binPath, "output", "O", "", "location to output binaries")
	viper.BindPFlag("main_path", buildCmd.Flags().Lookup("main_path"))
	viper.BindPFlag("index_path", buildCmd.Flags().Lookup("index_path"))
	viper.BindPFlag("bin_path", buildCmd.Flags().Lookup("output"))
	rootCmd.AddCommand(buildCmd)
}

// buildCmd is the command for building the project binary.
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Golang binaries and Content Management UI deployment",
	Long:  `Build Golang binaries and Content Management UI deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfigs()

		if !viper.InConfig("main_path") {
			log.Fatal("configuration not set: main_path")
		}
		if !viper.InConfig("bin_path") {
			log.Fatal("configuration not set: bin_path")
		}
		mainPath := viper.GetString("main_path")
		exists, err := pathExists(mainPath)
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			log.Fatal("main path does not exist")
		}

		binPath := viper.GetString("bin_path")
		if err := runCommand(fmt.Sprintf("go build -o %s %s", binPath, mainPath)); err != nil {
			log.Fatal(err)
		}
	},
}
