package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
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

func generateMain(modfile string, modules []string) error {
	if err := generateDir("bin", genHeadOnly, skipIfExists); err != nil {
		return err
	}

	ctx := struct {
		Modfile string
		Modules []string
	}{
		Modfile: modfile,
		Modules: modules,
	}

	if err := generateFile("main.go", MainTmpl, ctx, deleteIfExists); err != nil {
		return err
	}

	if err := RunCommand("go", "build", "-o", "bin/cosys", "main.go"); err != nil {
		return err
	}

	if err := os.Remove("main.go"); err != nil {
		return err
	}

	return nil
}

var MainTmpl = `package main

{{$modfile := .Modfile}}
import (
	"log"	

	"github.com/cosys-io/cosys/common"
{{range .Modules}}	"{{$modfile}}/modules/{{.}}"
{{end}})

func main() {
	var err error 
	
	modules := map[string]*common.Module{
{{range .Modules}}		"{{.}}": {{.}}.Module,
{{end}}}

	cosys := common.NewCosys(nil)

	cosys, err = cosys.Register(modules)
	if err != nil {
		log.Fatal(err)
	}

	if err := cosys.Start(); err != nil {
		log.Fatal(err)
	}
}
`
