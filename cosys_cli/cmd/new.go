package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/cosys_cli/cmd/generator"
	"github.com/iancoleman/strcase"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	modFile string
	db      string
	tmpl    string
)

func init() {
	newCmd.Flags().StringVarP(&modFile, "modfile", "M", "", "module name for go.mod file")
	newCmd.MarkFlagRequired("modfile")
	newCmd.Flags().StringVarP(&db, "database", "D", "sqlite3", "database system for the new project")
	newCmd.Flags().StringVarP(&tmpl, "template", "T", "", "template for the new project")

	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new project_name -P profile_name",
	Short: "Create a new cosys project",
	Long:  `Create a new cosys project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strcase.ToSnake(args[0])

		err := generateProject(projectName, modFile, db, tmpl)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func generateProject(projectName, modFile, db, tmpl string) error {
	var modules []string
	switch tmpl {
	default:
		modules = []string{
			"github.com/cosys-io/cosys/modules/" + db,
			"github.com/cosys-io/cosys/modules/server",
			"github.com/cosys-io/cosys/modules/admin",
		}
	}

	if err := generateConfigs(projectName); err != nil {
		return err
	}

	if err := generateModFile(projectName, modFile); err != nil {
		return err
	}

	if err := generateMain(projectName, modules); err != nil {
		return err
	}

	//if err := installModules(modules, Dir(projectName)); err != nil {
	//	return err
	//}

	return nil
}

func generateConfigs(projectName string) error {
	ctx := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	if err := gen.NewFile(filepath.Join(projectName, ".env"), EnvTmpl, ctx).Act(); err != nil {
		return err
	}

	if err := gen.NewFile(filepath.Join(projectName, ".cli_configs"), CliConfigTmpl, ctx).Act(); err != nil {
		return err
	}

	return nil
}

func generateModFile(projectName, modFile string) error {
	if err := RunCommand(fmt.Sprintf("go mod init %s", modFile), Dir(projectName)); err != nil {
		return err
	}

	//if err := RunCommand("go mod tidy", Dir(projectName)); err != nil {
	//	return err
	//}

	return nil
}

func generateMain(projectName string, modules []string) error {
	generator := gen.NewGenerator(
		gen.NewDir(filepath.Join(projectName, "cmd"), gen.GenHeadOnly),
		gen.NewDir(filepath.Join(projectName, "cmd", "cosys"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(projectName, "cmd", "cosys", "main.go"), MainTmpl, modules),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

var EnvTmpl = `HOST = localhost
PORT = 3000

DBNAME = cosys
DBHOST = db
DBPORT = 3000
DBUSER = cosys
DBPASS = cosys`

var CliConfigTmpl = `main_path: cmd/{{.ProjectName}}/main.go
index_path: web/
bin_path: bin/{{.ProjectName}}
`

var MainTmpl = `package main

import (
	"log"
{{range .}}    _ "{{.}}"
{{end}}
	"github.com/cosys-io/cosys/common"
)

func main() {
	cosys := common.NewCosys()

	if err := cosys.Start(); err != nil {
		log.Fatal(err)
	}
}`
