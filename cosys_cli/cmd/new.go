package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/cosys_cli/cmd/generator"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	profile string
	db      string
	tmpl    string
)

func init() {
	newCmd.Flags().StringVarP(&profile, "profile", "P", "", "profile name for the new project go module")
	newCmd.MarkFlagRequired("profile")
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
		projectName := args[0]

		err := generateProject(projectName, profile, db, tmpl)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func generateProject(projectName, pf, db, tmpl string) error {
	if err := generateConfigs(projectName, db, tmpl); err != nil {
		return err
	}

	if err := generateModules(projectName, pf, db, tmpl); err != nil {
		return err
	}

	if err := generateModfile(projectName, pf); err != nil {
		return err
	}

	if err := installModules([]string{db, "server", "module_service", "admin"}, Dir(projectName), Quiet); err != nil {
		return err
	}

	return nil
}

func generateConfigs(projectName, db, tmpl string) error {
	configsDir := filepath.Join(projectName, "configs")

	ctx := struct {
		Database string
	}{
		Database: db,
	}

	generator := gen.NewGenerator(
		gen.NewFile(filepath.Join(projectName, ".env"), EnvTmpl, ctx),
		gen.NewDir(configsDir),
		gen.NewFile(filepath.Join(configsDir, "admin.yaml"), AdminCfgTmpl, ctx),
		gen.NewFile(filepath.Join(configsDir, "database.yaml"), DbCfgTmpl, ctx),
		gen.NewFile(filepath.Join(configsDir, "module.yaml"), ModuleCfgTmpl, ctx),
		gen.NewFile(filepath.Join(configsDir, "server.yaml"), ServerCfgTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

func generateModules(projectName, pf, db, tmpl string) error {
	modulesDir := filepath.Join(projectName, "modules")

	if err := gen.NewDir(modulesDir).Act(); err != nil {
		return err
	}

	modfile := fmt.Sprintf("github.com/%s/%s", pf, projectName)

	if err := generateModule(filepath.Join(modulesDir, "api"), "api", modfile); err != nil {
		return err
	}

	return nil
}

func generateModfile(projectName, pf string) error {
	if err := RunCommand(fmt.Sprintf("go mod init github.com/%s/%s", pf, projectName), Dir(projectName)); err != nil {
		return err
	}

	//if err := RunCommand("go mod tidy", Dir(projectName)); err != nil {
	//	return err
	//}

	return nil
}

var EnvTmpl = `HOST = localhost
PORT = 3000

DBNAME = cosys
DBHOST = db
DBPORT = 3000
DBUSER = cosys
DBPASS = cosys`

var AdminCfgTmpl = ``

var DbCfgTmpl = `client: {{.Database}}
name: ENV.DBNAME
host: ENV.DBHOST
port: ENV.DBPORT
user: ENV.DBUSER
pass: ENV.DBPASS`

var ModuleCfgTmpl = `modules: 
  - api
  - module_service
  - server
  - admin
  - sqlite3`

var ServerCfgTmpl = `host: ENV.HOST
port: ENV.PORT`
