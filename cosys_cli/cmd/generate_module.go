package cmd

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(generateModuleCmd)
}

var generateModuleCmd = &cobra.Command{
	Use:   "module module_name",
	Short: "Generate an module",
	Long:  "Generate an module.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]

		modfile, err := getModfile()
		if err != nil {
			log.Fatal(err)
		}

		if err := generateModule(filepath.Join("modules/", moduleName), moduleName, modfile); err != nil {
			log.Fatal(err)
		}
	},
}

func generateModule(moduleDir, moduleName, modfile string) error {
	if err := generateDir(moduleDir, genHeadOnly); err != nil {
		return err
	}

	ctx := struct {
		ModuleName string
		Modfile    string
	}{
		moduleName,
		modfile,
	}

	if err := generateFile(filepath.Join(moduleDir, "module.go"), ModuleTmpl, ctx); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "content_types"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "content_types", "models.go"), ModelsTmpl, nil); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "controllers"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "controllers", "controllers.go"), ControllersTmpl, nil); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "middlewares"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "middlewares", "middlewares.go"), MiddlewaresTmpl, nil); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "policies"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "policies", "policies.go"), PoliciesTmpl, nil); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "routes"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "routes", "routes.go"), RoutesTmpl, nil); err != nil {
		return err
	}

	if err := generateDir(filepath.Join(moduleDir, "services"), genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(moduleDir, "services", "services.go"), ServicesTmpl, nil); err != nil {
		return err
	}

	return nil
}
