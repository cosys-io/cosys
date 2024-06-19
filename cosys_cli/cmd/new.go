package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/cosys_cli/cmd/generator"
	"log"
	"os/exec"
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

	return nil
}

func generateConfigs(projectName, db, tmpl string) error {
	configsDir := filepath.Join(projectName, "configs")

	generator := gen.NewGenerator(
		gen.NewDir(configsDir),
		gen.NewFile(filepath.Join(configsDir, "admin.go"), AdminCfgTmpl, nil),
		gen.NewFile(filepath.Join(configsDir, "database.go"), DbCfgTmpl, nil),
		gen.NewFile(filepath.Join(configsDir, "server.go"), ServerCfgTmpl, nil),
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
	cmd := exec.Command("go", "mod", "init", fmt.Sprintf("github.com/%s/%s", pf, projectName))
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "get", "github.com/cosys-io/cosys")
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

var AdminCfgTmpl = `package configs`

var DbCfgTmpl = `package configs`

var ServerCfgTmpl = `package configs`
