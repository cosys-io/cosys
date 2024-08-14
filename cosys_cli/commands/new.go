package commands

import (
	"fmt"
	gen "github.com/cosys-io/cosys/cosys_cli/generator"
	"github.com/iancoleman/strcase"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	modFile string // modFile is bound to the modfile flag, which is required.
	db      string // db is bound to the database flag.
	tmpl    string // tmpl is bound to the template flag.
)

func init() {
	newCmd.Flags().StringVarP(&modFile, "modfile", "M", "", "module name for go.mod file")
	newCmd.MarkFlagRequired("modfile")
	newCmd.Flags().StringVarP(&db, "database", "D", "sqlite3", "database system for the new project")
	newCmd.Flags().StringVarP(&tmpl, "template", "T", "", "template for the new project")

	rootCmd.AddCommand(newCmd)
}

// newCmd is the command for creating a new cosys project.
var newCmd = &cobra.Command{
	Use:   "new project_name -M module_name",
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

// generateProject generates files for a new project.
func generateProject(projectName, modFile, db, tmpl string) error {
	var modules []string
	switch tmpl {
	default:
		modules = []string{
			"github.com/cosys-io/cosys/modules/" + db,
			"github.com/cosys-io/cosys/modules/server",
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

// generateConfigs generates files for project configuration.
func generateConfigs(projectName string) error {
	ctx := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	if err := gen.NewFile(filepath.Join(projectName, ".env"), envTmpl, ctx).Act(); err != nil {
		return err
	}

	if err := gen.NewFile(filepath.Join(projectName, ".cli_configs"), cliConfigTmpl, ctx).Act(); err != nil {
		return err
	}

	return nil
}

// generateModFile creates a new go module.
func generateModFile(projectName, modFile string) error {
	if err := runCommand(fmt.Sprintf("go mod init %s", modFile), dir(projectName)); err != nil {
		return err
	}

	//if err := runCommand("go mod tidy", Dir(projectName)); err != nil {
	//	return err
	//}

	return nil
}

// generateMain generates the main package.
func generateMain(projectName string, modules []string) error {
	generator := gen.NewGenerator(
		gen.NewDir(filepath.Join(projectName, "cmd"), gen.GenHeadOnly),
		gen.NewDir(filepath.Join(projectName, "cmd", "cosys"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(projectName, "cmd", "cosys", "main.go"), mainTmpl, modules),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// envTmpl is the template for the .env file.
var envTmpl = `HOST = localhost
PORT = 3000

DBNAME = cosys
DBHOST = db
DBPORT = 3000
DBUSER = cosys
DBPASS = cosys`

// cliConfigTmpl is the template for the .cli_config file.
var cliConfigTmpl = `main_path: cmd/{{.ProjectName}}/main.go
index_path: web/
bin_path: bin/{{.ProjectName}}
`

// mainTmpl is the template for the main file.
var mainTmpl = `package main

import (
	"log"
{{range .}}    _ "{{.}}"
{{end}}
	"github.com/cosys-io/cosys/common"
)

func main() {
	cosys, err := common.New()
	if err != nil {
		log.Fatal(err)
	}

	if err = cosys.Start(); err != nil {
		log.Fatal(err)
	}
}`
