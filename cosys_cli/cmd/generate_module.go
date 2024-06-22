package cmd

import (
	gen "github.com/cosys-io/cosys/cosys_cli/cmd/generator"
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

		modFile, err := getModFile()
		if err != nil {
			log.Fatal(err)
		}

		if err := generateModule(filepath.Join("modules/", moduleName), moduleName, modFile); err != nil {
			log.Fatal(err)
		}
	},
}

func generateModule(moduleDir, moduleName, modFile string) error {
	ctx := struct {
		ModuleName string
		ModFile    string
	}{
		moduleName,
		modFile,
	}

	generator := gen.NewGenerator(
		gen.NewDir(moduleDir, gen.GenHeadOnly),
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
		gen.NewDir(filepath.Join(moduleDir, "services"), gen.GenHeadOnly),
		gen.NewFile(filepath.Join(moduleDir, "services", "services.go"), ServicesTmpl, nil),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

var ModuleTmpl = `package {{.ModuleName}}

import (
	"github.com/cosys-io/cosys/common"
	"{{.ModFile}}/modules/{{.ModuleName}}/controllers"
	"{{.ModFile}}/modules/{{.ModuleName}}/middlewares"
	"{{.ModFile}}/modules/{{.ModuleName}}/policies"
	"{{.ModFile}}/modules/{{.ModuleName}}/routes"
	"{{.ModFile}}/modules/{{.ModuleName}}/content_types"
)

var Module = &common.Module{
	Routes: routes.Routes,
	Controllers: controllers.Controllers,
	Middlewares: middlewares.Middlewares,
	Policies: policies.Policies,
	Models: models.Models,
	Services: nil,

	OnRegister: nil,
	OnDestroy: nil,
}
`

var ControllersTmpl = `package controllers

import "github.com/cosys-io/cosys/common"

var Controllers = map[string]common.Controller{
}`

var MiddlewaresTmpl = `package middlewares

import "github.com/cosys-io/cosys/common"

var Middlewares = map[string]common.Middleware{
}`

var PoliciesTmpl = `package policies

import "github.com/cosys-io/cosys/common"

var Policies = map[string]common.Policy{
}`

var RoutesTmpl = `package routes

import "github.com/cosys-io/cosys/common"

var Routes = []*common.Route{
}`

var ServicesTmpl = `package services

import "github.com/cosys-io/cosys/common"

var Services = map[string]common.Service{
}`

var ModelsTmpl = `package models

import (
	"github.com/cosys-io/cosys/common"
)

var Models = map[string]common.Model{
}`
