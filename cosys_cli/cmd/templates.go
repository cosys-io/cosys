package cmd

var ModTmpl = `module github.com/{{.Profile}}/{{.ProjectName}}

go 1.22.2`

var AdminCfgTmpl = `package configs`

var DbCfgTmpl = `package configs`

var ServerCfgTmpl = `package configs`

var ModuleTmpl = `package {{.ModuleName}}

import (
	"github.com/cosys-io/cosys/common"
	"{{.Modfile}}/modules/{{.ModuleName}}/controllers"
	"{{.Modfile}}/modules/{{.ModuleName}}/middlewares"
	"{{.Modfile}}/modules/{{.ModuleName}}/policies"
	"{{.Modfile}}/modules/{{.ModuleName}}/routes"
)

var Module = &common.Module{
	Routes: routes.Routes,
	Controllers: controllers.Controllers,
	Middlewares: middlewares.Middlewares,
	Policies: policies.Policies,
	Models: nil,
	Services: nil,

	OnRegister: nil,
	OnDestroy: nil,
}
`

var ControllersTmpl = `package controllers

import "github.com/cosys-io/cosys/common"

var Controllers = map[string]*common.Controller{}`

var MiddlewaresTmpl = `package middlewares

import "github.com/cosys-io/cosys/common"

var Middlewares = map[string]common.Middleware{}`

var PoliciesTmpl = `package policies

import "github.com/cosys-io/cosys/common"

var Policies = map[string]common.Policy{}`

var RoutesTmpl = `package routes

import "github.com/cosys-io/cosys/common"

var Routes = []*common.Route{}`

var ServicesTmpl = `package services

import "github.com/cosys-io/cosys/common"

var Services = map[string]*common.Service{}`

var ModelsTmpl = `package models

import "github.com/cosys-io/cosys/common"

var Models = map[string]*common.Model{}`

var ModelTmpl = `package {{.CollectionName}}
	
import "github.com/cosys-io/cosys/common"

type {{.SingularName}} struct {
{{range .Attributes}}    {{.NamePascal}} {{.TypeLower}} ` + "`" + `json:"{{.NameCamel}}"` + "`" + `
{{end}}}


type {{.PluralName}}Model struct {
	schema    *common.ModelSchema
	lifecycle common.Lifecycle

{{range .Attributes}}    {{.NamePascal}} *common.{{.TypeUpper}}Attribute
{{end}}}

var {{.PluralName}} = {{.PluralName}}Model{
	Schema,
	Lifecycle,

{{range .Attributes}}    common.New{{.TypeUpper}}Attribute("{{.NameCamel}}", "{{.NamePascal}}"),
{{end}}}

func (m {{.PluralName}}Model) Name_() string {
	return "{{.CollectionName}}"
}

func (m {{.PluralName}}Model) New_() common.Entity {
	return &{{.SingularName}}{}
}

{{$Model := .PluralName}}
func (m {{.PluralName}}Model) All_() []common.Attribute {
	return []common.Attribute{
{{range .Attributes}}        {{$Model}}.{{.NamePascal}},
{{end}}}
}

func (m {{.PluralName}}Model) Id_() *common.IntAttribute {
	return {{.PluralName}}.Id
}

func (m {{.PluralName}}Model) Schema_() *common.ModelSchema {
	return m.schema
}

func (m {{.PluralName}}Model) Lifecycle_() common.Lifecycle {
	return m.lifecycle
}`

var LifecycleTmpl = `package {{.CollectionName}}

import "github.com/cosys-io/cosys/common"

var Lifecycle = common.NewLifeCycle()`

var SchemaGoTmpl = `package {{.CollectionName}}

import "github.com/cosys-io/cosys/common"

var Schema = &common.ModelSchema{
	ModelType:      "collectiontype",
	CollectionName: "{{.CollectionName}}",
	DisplayName:    "{{.DisplayName}}",
	SingularName:   "{{.SingularName}}",
	PluralName:     "{{.PluralName}}",
	Description:    "{{.Description}}",
	Attributes: map[string]*common.AttributeSchema{
{{range $name, $attr := .Attributes}}        "{{$name}}": {
			Type: "{{$attr.Type}}",
			
			Required: {{$attr.Required}},
			Max: {{$attr.Max}},
			Min: {{$attr.Min}},
			MaxLength: {{$attr.MaxLength}},
			MinLength: {{$attr.MinLength}},
			Private: {{$attr.Private}},
			NotConfigurable: {{$attr.NotConfigurable}},
			
			Default: "{{$attr.Default}}",
			NotNullable: {{$attr.NotNullable}},
			Unsigned: {{$attr.Unsigned}},
			Unique: {{$attr.Unique}},
		},
{{end}}    },
}`

var SchemaYamlTmpl = `modelType: {{.ModelType}}
collectionName: {{.CollectionName}}
displayName: {{.DisplayName}}
singularName: {{.SingularName}}
pluralName: {{.PluralName}}
description: {{.Description}}
attributes:
{{range $name, $attr := .Attributes}}  {{$name}}:
		type: {{$attr.Type}}{{if $attr.Required}}
		required: true{{end}}{{if ne $attr.Max 2147483647}}
		max: {{.Max}}{{end}}{{if ne $attr.Min -2147483648}}
		min: {{.Min}}{{end}}{{if ne $attr.MaxLength -1}}
		maxLength: {{.MaxLength}}{{end}}{{if ne $attr.MinLength -1}}
		minLength: {{.MinLength}}{{end}}{{if $attr.Private}}
		private: true{{end}}{{if $attr.NotConfigurable}}
		notConfigurable: true{{end}}{{if $attr.Default}}
		default: {{.Default}}{{end}}{{if $attr.NotNullable}}
		notNullable: true{{end}}{{if $attr.Unsigned}}
		unsigned: true{{end}}{{if $attr.Unique}}
		unique: true{{end}}
{{end}}`

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
