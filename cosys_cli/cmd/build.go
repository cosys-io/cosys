package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	main_path  string
	index_path string
	bin_path   string
)

func init() {
	buildCmd.Flags().StringVarP(&main_path, "main_path", "M", "", "location of main package")
	buildCmd.Flags().StringVarP(&index_path, "index_path", "I", "", "location of ui index file")
	buildCmd.Flags().StringVarP(&bin_path, "output", "O", "", "location to output binaries")
	viper.BindPFlag("main_path", buildCmd.Flags().Lookup("main_path"))
	viper.BindPFlag("index_path", buildCmd.Flags().Lookup("index_path"))
	viper.BindPFlag("bin_path", buildCmd.Flags().Lookup("output"))
	rootCmd.AddCommand(buildCmd)
}

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
		if err := RunCommand(fmt.Sprintf("go build -o %s %s", binPath, mainPath)); err != nil {
			log.Fatal(err)
		}
	},
}
