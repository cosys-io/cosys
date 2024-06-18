package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	displayName  string
	singularName string
	pluralName   string
	description  string
)

func init() {
	generateCollectionCmd.Flags().StringVarP(&displayName, "display", "N", "", "display name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&singularName, "singular", "S", "", "singular name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&pluralName, "plural", "P", "", "plural name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&description, "description", "D", "", "description of the new content type")
	generateCollectionCmd.MarkFlagRequired("display")
	generateCollectionCmd.MarkFlagRequired("singular")
	generateCollectionCmd.MarkFlagRequired("plural")

	generateCmd.AddCommand(generateCollectionCmd)
}

var generateCollectionCmd = &cobra.Command{
	Use:   "collection collection_name [attributes] [flags]",
	Short: "Generate a collection type",
	Long:  "Generate a collection type.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collectionName := args[0]

		attributes := args[1:]

		schema, ctx, err := schemaFromArgs(collectionName, displayName, singularName, pluralName, description, attributes)
		if err != nil {
			log.Fatal(err)
		}

		if err := generateType(collectionName, schema, ctx); err != nil {
			log.Fatal(err)
		}
	},
}

func generateType(typeName string, schema *common.ModelSchema, ctx *ModelCtx) error {
	if err := generateModel(typeName, schema, ctx); err != nil {
		return nil
	}

	if err := generateApi(typeName, ctx); err != nil {
		return nil
	}

	return nil
}

func generateModel(typeName string, schema *common.ModelSchema, ctx *ModelCtx) error {
	typeDir := filepath.Join("modules/api/content_types", typeName)

	if err := generateDir(typeDir, genHeadOnly); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "schema.yaml"), SchemaYamlTmpl, schema); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "schema.go"), SchemaGoTmpl, schema); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "lifecycle.go"), LifecycleTmpl, ctx); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "model.go"), ModelTmpl, ctx); err != nil {
		return err
	}

	return nil
}

func generateApi(typeName string, ctx *ModelCtx) error {
	if err := generateFile(filepath.Join("modules/api/controllers", typeName+"_controllers.go"), ModelControllerTmpl, ctx); err != nil {
		return err
	}

	if err := modifyFile("modules/api/content_types/models.go", `import \(`, ModelsImportTmpl, ctx); err != nil {
		return err
	}

	if err := modifyFile("modules/api/content_types/models.go", `var Models = map\[string\]common\.Model\{`, ModelsStructTmpl, ctx); err != nil {
		return err
	}

	if err := modifyFile("modules/api/controllers/controllers.go", `var Controllers = map\[string\]common\.Controller\{`, ControllersStructTmpl, ctx); err != nil {
		return err
	}

	if err := modifyFile("modules/api/routes/routes.go", `var Routes = \[\]\*common\.Route\{`, RoutesStructTmpl, ctx); err != nil {
		return err
	}

	return nil
}

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

var ModelControllerTmpl = `package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosys-io/cosys/common"
)

var {{.PluralName}}Controller = map[string]common.Action{
	"findOne": findOne{{.SingularName}},
	"create":  create{{.SingularName}},
	"update":  update{{.SingularName}},
	"delete":  delete{{.SingularName}},
}

func findOne{{.SingularName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
			return
		}

		id, err := strconv.Atoi(params[0])
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		entity, err := cs.ModuleService().FindOne("api.{{.CollectionName}}", id, common.MSParam())
		if err != nil {
			common.RespondError(w, "Could not find {{.SingularNameCamel}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, entity, http.StatusOK)
	}
}

func create{{.SingularName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model, ok := cs.Models["api.{{.CollectionName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		entity := model.New_()

		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		newEntity, err := cs.ModuleService().Create("api.{{.CollectionName}}", entity, common.MSParam())
		if err != nil {
			common.RespondError(w, "Could not create {{.SingularNameCamel}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, newEntity, http.StatusOK)
	}
}

func update{{.SingularName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
			return
		}

		id, err := strconv.Atoi(params[0])
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		model, ok := cs.Models["api.{{.CollectionName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		entity := model.New_()

		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		newEntity, err := cs.ModuleService().Update("api.{{.CollectionName}}", entity, id, common.MSParam().SetField(model.All_()...))
		if err != nil {
			common.RespondError(w, "Could not update {{.SingularNameCamel}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, newEntity, http.StatusOK)
	}
}

func delete{{.SingularName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
			return
		}

		id, err := strconv.Atoi(params[0])
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		oldEntity, err := cs.ModuleService().Delete("api.{{.CollectionName}}", id, common.MSParam())
		if err != nil {
			common.RespondError(w, "Could not delete {{.SingularNameCamel}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, oldEntity, http.StatusOK)
	}
}`

var ModelsImportTmpl = `import (
	"{{.ModFile}}/modules/api/content_types/{{.CollectionName}}"`

var ModelsStructTmpl = `var Models = map[string]common.Model{
	"api.{{.CollectionName}}": {{.CollectionName}}.{{.PluralName}},`

var ControllersStructTmpl = `var Controllers = map[string]common.Controller{
	"{{.CollectionName}}": {{.PluralName}}Controller,`

var RoutesStructTmpl = `var Routes = []*common.Route{
	common.NewRoute("GET", ` + "`/{{.CollectionName}}/([0-9]+)`" + `, "{{.CollectionName}}.findOne"),
	common.NewRoute("POST", ` + "`/{{.CollectionName}}`" + `, "{{.CollectionName}}.create"),
	common.NewRoute("PUT", ` + "`/{{.CollectionName}}/([0-9]+)`" + `, "{{.CollectionName}}.update"),
	common.NewRoute("DELETE", ` + "`/{{.CollectionName}}/([0-9]+)`" + `, "{{.CollectionName}}.delete"),`

func schemaFromArgs(collection string, display string, singular string, plural string, description string, attrStrings []string) (*common.ModelSchema, *ModelCtx, error) {
	modelSchema := &common.ModelSchema{
		ModelType:      "collectionType",
		CollectionName: collection,
		DisplayName:    display,
		SingularName:   singular,
		PluralName:     plural,
		Description:    description,
		Attributes:     map[string]*common.AttributeSchema{},
	}

	caser := cases.Title(language.English)

	modfile, err := getModfile()
	if err != nil {
		return nil, nil, err
	}

	modelCtx := &ModelCtx{
		ModFile:           modfile,
		CollectionName:    collection,
		SingularName:      caser.String(singular),
		SingularNameCamel: singular,
		PluralName:        caser.String(plural),
		Attributes: []*AttributeCtx{
			{
				TypeLower:  "int",
				TypeUpper:  "Int",
				NameCamel:  "id",
				NamePascal: "Id",
			},
		},
	}

	for _, attrString := range attrStrings {
		split := strings.Split(attrString, ":")
		if len(split) < 2 {
			return nil, nil, fmt.Errorf("invalid attribute format: %s", attrString)
		}
		attrName := split[0]
		attrType := split[1]

		attrSchema := &common.AttributeSchema{
			Type:            attrType,
			Required:        false,
			Max:             2147483647,
			Min:             -2147483648,
			MaxLength:       -1,
			MinLength:       -1,
			Private:         false,
			NotConfigurable: false,
			Default:         "",
			NotNullable:     false,
			Unsigned:        false,
			Unique:          false,
		}

		for _, option := range split[2:] {
			switch {
			case option == "required":
				attrSchema.Required = true
			case regexp.MustCompile(`^max=([-0-9]+)$`).MatchString(option):
				matches := regexp.MustCompile(`^max=([-0-9]+)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				attrSchema.Max = val
			case regexp.MustCompile(`^min=([-0-9]+)$`).MatchString(option):
				matches := regexp.MustCompile(`^min=([-0-9]+)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				attrSchema.Min = val
			case regexp.MustCompile(`^maxlength=(\d+)$`).MatchString(option):
				matches := regexp.MustCompile(`^maxlength=(\d+)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, nil, err
				}
				attrSchema.MaxLength = val
			case regexp.MustCompile(`^minlength=(\d+)$`).MatchString(option):
				matches := regexp.MustCompile(`^minlength=(\d+)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, nil, err
				}
				attrSchema.MinLength = val
			case option == "private":
				attrSchema.Private = true
			case option == "notconfigurable":
				attrSchema.NotConfigurable = true
			case regexp.MustCompile(`^default=(.+)=$`).MatchString(option):
				matches := regexp.MustCompile(`^default=(.+)$`).FindStringSubmatch(option)
				attrSchema.Default = matches[1]
			case option == "notnullable":
				attrSchema.NotNullable = true
			case option == "unsigned":
				attrSchema.Unsigned = true
			case option == "unique":
				attrSchema.Unique = true
			}
		}

		if _, ok := modelSchema.Attributes[attrName]; ok {
			return nil, nil, fmt.Errorf("duplicate attribute: %s", attrName)
		}

		attrCtx := &AttributeCtx{
			NameCamel:  attrName,
			NamePascal: caser.String(attrName),
		}

		switch attrType {
		case "integer":
			attrCtx.TypeLower = "int"
			attrCtx.TypeUpper = "Int"
		case "string":
			attrCtx.TypeLower = "string"
			attrCtx.TypeUpper = "String"
		case "boolean":
			attrCtx.TypeLower = "bool"
			attrCtx.TypeUpper = "Bool"
		}

		modelSchema.Attributes[attrName] = attrSchema
		modelCtx.Attributes = append(modelCtx.Attributes, attrCtx)
	}

	return modelSchema, modelCtx, nil
}

type ModelCtx struct {
	ModFile           string
	CollectionName    string
	SingularName      string
	SingularNameCamel string
	PluralName        string
	Attributes        []*AttributeCtx
}

type AttributeCtx struct {
	TypeLower  string
	TypeUpper  string
	NameCamel  string
	NamePascal string
}
