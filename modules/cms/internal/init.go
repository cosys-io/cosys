package internal

import (
	"github.com/cosys-io/cosys/common"
	gen "github.com/cosys-io/cosys/cosys_cli/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init module_path",
	Short: "Generate default configurations and configurations for the cms module",
	Long:  `Generate default configurations and configurations for the cms module.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modulePath := args[0]

		modFile, err := getModFile()
		if err != nil {
			log.Fatal(err)
		}

		common.InitConfigs()
		viper.Set("cms_content_types_path", filepath.Join(modulePath, "content_types"))
		viper.Set("cms_routes_path", filepath.Join(modulePath, "routes"))
		viper.Set("cms_controllers_path", filepath.Join(modulePath, "controllers"))
		viper.Set("cms_middlewares_path", filepath.Join(modulePath, "middlewares"))
		viper.Set("cms_policies_path", filepath.Join(modulePath, "policies"))
		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}

		if err := generateModule(modulePath, "cms", modFile); err != nil {
			log.Fatal(err)
		}

	},
}

func generateModule(moduleDir, moduleName, modFile string) error {
	ctx := struct {
		ModuleDir  string
		ModuleName string
		ModFile    string
	}{
		moduleDir,
		moduleName,
		modFile,
	}

	generator := gen.NewGenerator(
		gen.NewDir(moduleDir),
		gen.NewFile(filepath.Join(moduleDir, "module.go"), ModuleTmpl, ctx),
		gen.NewDir(filepath.Join(moduleDir, "content_types"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "content_types", "models.go"), ModelsTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "controllers"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "controllers", "controllers.go"), ControllersTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "middlewares"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "middlewares", "middlewares.go"), MiddlewaresTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "policies"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "policies", "policies.go"), PoliciesTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "routes"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "routes", "routes.go"), RoutesTmpl, nil),
		//gen.NewDir(filepath.Join(moduleDir, "services"), gen.GenHeadOnly),
		//gen.NewFile(filepath.Join(moduleDir, "services", "services.go"), ServicesTmpl, nil),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

var ModuleTmpl = `package {{.ModuleName}}

import (
	"github.com/cosys-io/cosys/common"
	"{{.ModFile}}/{{.ModuleDir}}/content_types"
	"{{.ModFile}}/{{.ModuleDir}}/controllers"
	"{{.ModFile}}/{{.ModuleDir}}/middlewares"
	"{{.ModFile}}/{{.ModuleDir}}/policies"
	"{{.ModFile}}/{{.ModuleDir}}/routes"
	"github.com/cosys-io/cosys/modules/cms/admin"
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		if err := cosys.AddRoutes(routes.Routes...); err != nil {
			return err
		}

		if err := cosys.AddControllers(controllers.Controllers...); err != nil {
			return err
		}

		if err := cosys.AddMiddlewares(middlewares.Middlewares...); err != nil {
			return err
		}

		if err := cosys.AddPolicies(policies.Policies...); err != nil {
			return err
		}

		if err := cosys.AddModels(models.Models); err != nil {
			return err
		}

		if err := admin.AddAdminRoutes(cosys, models.Models); err != nil {
			return err
		}

		if err := admin.AddSchemaRoutes(cosys, models.Models); err != nil {
			return err
		}

		return nil
	})
}
`

var ControllersTmpl = `package controllers

import "github.com/cosys-io/cosys/common"

var Controllers = []common.Controller{
}`

var MiddlewaresTmpl = `package middlewares

import "github.com/cosys-io/cosys/common"

var Middlewares = []common.Middleware{
}`

var PoliciesTmpl = `package policies

import "github.com/cosys-io/cosys/common"

var Policies = []common.Policy{
}`

var RoutesTmpl = `package routes

import "github.com/cosys-io/cosys/common"

var Routes = []common.Route{
}`

var ModelsTmpl = `package models

import (
	"github.com/cosys-io/cosys/common"
)

var Models = map[string]common.Model{
}`
