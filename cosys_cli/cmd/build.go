package cmd

import (
	"bufio"
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Golang binaries and Content Management UI deployment",
	Long:  `Build Golang binaries and Content Management UI deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		modfile, err := getModfile()
		if err != nil {
			log.Fatal(err)
		}

		modules, err := getModules()
		if err != nil {
			log.Fatal(err)
		}

		if err := generateMain(modfile, modules); err != nil {
			log.Fatal(err)
		}
	},
}

func getModules() ([]string, error) {
	entries, err := os.ReadDir("./modules")
	if err != nil {
		return nil, err
	}

	var modules []string

	for _, entry := range entries {
		if entry.IsDir() {
			modules = append(modules, entry.Name())
		}
	}

	return modules, nil
}

func getModfile() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		return line[7:], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("module not found")
}

func generateMain(modfile string, modules []string) error {
	if err := generateDir("bin", genHeadOnly); err != nil {
		return nil
	}

	ctx := struct {
		Modfile string
		Modules []string
	}{
		Modfile: modfile,
		Modules: modules,
	}

	if err := generateFile("main.go", MainTmpl, ctx, deleteIfExists); err != nil {
		return nil
	}

	cmd := exec.Command("go", "build", "-o", "bin/cosys", "main.go")
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Remove("main.go"); err != nil {
		return err
	}

	return nil
}
