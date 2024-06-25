package cmd

import (
	"github.com/cosys-io/cosys/cosys_cli/cmd/generator"
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
		modFile, err := getModFile()
		if err != nil {
			log.Fatal(err)
		}

		modules, err := getModules()
		if err != nil {
			log.Fatal(err)
		}

		if err := generateMain(modFile, modules); err != nil {
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

func generateMain(modFile string, modules []string) error {
	if err := gen.NewDir("bin", gen.GenHeadOnly, gen.SkipIfExists).Act(); err != nil {
		return err
	}

	ctx := struct {
		ModFile string
		Modules []string
	}{
		ModFile: modFile,
		Modules: modules,
	}

	if err := gen.NewFile("main.go", MainTmpl, ctx, gen.DeleteIfExists).Act(); err != nil {
		return err
	}

	defer os.Remove("main.go")

	if err := RunCommand("go build -o bin/cosys main.go"); err != nil {
		return err
	}

	return nil
}

var MainTmpl = `package main

{{$modFile := .ModFile}}
import (
	"log"	

	"github.com/cosys-io/cosys/common"
{{range .Modules}}	"{{$modFile}}/modules/{{.}}"
{{end}})

func main() {
	var err error 
	
	modules := map[string]*common.Module{
{{range .Modules}}		"{{.}}": {{.}}.Module,
{{end}}}

	cfg, err := common.GetConfigs("configs")
	if err != nil {
		log.Fatal(err)	
	}

	cosys := common.NewCosys(cfg)

	cosys, err = cosys.Register(modules)
	if err != nil {
		log.Fatal(err)
	}

	if err := cosys.Start(); err != nil {
		log.Fatal(err)
	}
}
`
