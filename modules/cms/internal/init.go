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

// initCmd is the command for generating the default configurations and code for the cms module.
var initCmd = &cobra.Command{
	Use:   "init module_path",
	Short: "Generate default configurations and code for the cms module",
	Long:  `Generate default configurations and code for the cms module.`,
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

// generateModule generates the code for the content types, controllers, middlewares, policies and routes packages.
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
		gen.NewFile(filepath.Join(moduleDir, "module.go"), moduleTmpl, ctx),
		gen.NewDir(filepath.Join(moduleDir, "content_types"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "content_types", "models.go"), modelsTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "controllers"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "controllers", "controllers.go"), controllersTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "middlewares"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "middlewares", "middlewares.go"), middlewaresTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "policies"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "policies", "policies.go"), policiesTmpl, nil),
		gen.NewDir(filepath.Join(moduleDir, "routes"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "routes", "routes.go"), routesTmpl, nil),
		//gen.NewDir(filepath.Join(moduleDir, "services"), gen.GenHeadOnly),
		//gen.NewFile(filepath.Join(moduleDir, "services", "services.go"), ServicesTmpl, nil),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// moduleTmpl is the template for creating the module.go file.
var moduleTmpl = `package {{.ModuleName}}

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

// controllersTmpl is the template for creating the controllers.go file.
var controllersTmpl = `package controllers

import "github.com/cosys-io/cosys/common"

var Controllers = []common.Controller{
}`

// middlewaresTmpl is the template for creating the middlewares.go file.
var middlewaresTmpl = `package middlewares

import "github.com/cosys-io/cosys/common"

var Middlewares = []common.Middleware{
}`

// policiesTmpl is the template for creating the policies.go file.
var policiesTmpl = `package policies

import "github.com/cosys-io/cosys/common"

var Policies = []common.Policy{
}`

// routesTmpl is the template for creating the routes.go file.
var routesTmpl = `package routes

import "github.com/cosys-io/cosys/common"

var Routes = []common.Route{
}`

// modelsTmpl is the template for creating the models.go file.
var modelsTmpl = `package models

import (
	"github.com/cosys-io/cosys/common"
)

var Models = map[string]common.Model{
}`
