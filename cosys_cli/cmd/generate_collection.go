package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	gen "github.com/cosys-io/cosys/cosys_cli/cmd/generator"
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

		schema, err := schemaFromArgs(collectionName, displayName, singularName, pluralName, description, attributes)
		if err != nil {
			log.Fatal(err)
		}

		if err := GenerateType(collectionName, schema); err != nil {
			log.Fatal(err)
		}
	},
}

func GenerateType(typeName string, schema *common.ModelSchema) error {
	ctx, err := ctxFromSchema(schema)
	if err != nil {
		return err
	}

	if err = generateModel(typeName, schema, ctx); err != nil {
		return err
	}

	if err = generateApi(typeName, ctx); err != nil {
		return err
	}

	return nil
}

func generateModel(typeName string, schema *common.ModelSchema, ctx *ModelCtx) error {
	typeDir := filepath.Join("modules/api/content_types", typeName)

	generator := gen.NewGenerator(
		gen.NewDir(typeDir, gen.GenHeadOnly),
		gen.NewFile(filepath.Join(typeDir, "schema.yaml"), SchemaYamlTmpl, schema),
		gen.NewFile(filepath.Join(typeDir, "lifecycle.go"), LifecycleTmpl, ctx),
		gen.NewFile(filepath.Join(typeDir, "model.go"), ModelTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

func generateApi(typeName string, ctx *ModelCtx) error {
	generator := gen.NewGenerator(
		gen.NewFile(filepath.Join("modules/api/controllers", typeName+"_controllers.go"), ModelControllerTmpl, ctx),
		gen.ModifyFile("modules/api/content_types/models.go", `import \(`, ModelsImportTmpl, ctx),
		gen.ModifyFile("modules/api/content_types/models.go", `var Models = map\[string\]common\.Model\{`, ModelsStructTmpl, ctx),
		gen.ModifyFile("modules/api/controllers/controllers.go", `var Controllers = map\[string\]common\.Controller\{`, ControllersStructTmpl, ctx),
		gen.ModifyFile("modules/api/routes/routes.go", `var Routes = \[\]\*common\.Route\{`, RoutesStructTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

var ModelTmpl = `package {{.CollectionName}}
	
import (
	"log"
	"github.com/cosys-io/cosys/common"
)

var (
	Schema = &common.ModelSchema{}
)

func init() {
	var err error
	Schema, err = common.GetSchema("modules/api/content_types/{{.CollectionName}}/schema.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

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
	return Schema
}

func (m {{.PluralName}}Model) Lifecycle_() common.Lifecycle {
	return m.lifecycle
}`

var LifecycleTmpl = `package {{.CollectionName}}

import "github.com/cosys-io/cosys/common"

var Lifecycle = common.NewLifeCycle()`

var SchemaYamlTmpl = `modelType: {{.ModelType}}
collectionName: {{.CollectionName}}
displayName: {{.DisplayName}}
singularName: {{.SingularName}}
pluralName: {{.PluralName}}
description: {{.Description}}
attributes:
{{range .Attributes}}  - name: {{.Name}}
    simpleDataType: {{.SimpleType}}
    detailedDataType: {{.DetailedType}}{{if not .ShownInTable}}
    shownInTable: false{{end}}{{if .Required}}
    required: true{{end}}{{if ne .Max 2147483647}}
    max: {{.Max}}{{end}}{{if ne .Min -2147483648}}
    min: {{.Min}}{{end}}{{if ne .MaxLength -1}}
    maxLength: {{.MaxLength}}{{end}}{{if ne .MinLength -1}}
    minLength: {{.MinLength}}{{end}}{{if .Private}}
    private: true{{end}}{{if not .Editable}}
    editable: false{{end}}{{if .Default}}
    default: {{.Default}}{{end}}{{if not .Nullable}}
    nullable: false{{end}}{{if .Unique}}
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

		newEntity, err := cs.ModuleService().Update("api.{{.CollectionName}}", entity, id, common.MSParam())
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

func schemaFromArgs(collection string, display string, singular string, plural string, description string, attrStrings []string) (*common.ModelSchema, error) {
	modelSchema := &common.ModelSchema{
		ModelType:      "collectionType",
		CollectionName: collection,
		DisplayName:    display,
		SingularName:   singular,
		PluralName:     plural,
		Description:    description,
		Attributes: []*common.AttributeSchema{
			{
				Name:         "id",
				SimpleType:   "Number",
				DetailedType: "Int",
				ShownInTable: true,
				Required:     true,
				Max:          2147483647,
				Min:          -2147483648,
				MaxLength:    -1,
				MinLength:    -1,
				Private:      false,
				Editable:     false,
				Nullable:     false,
				Unique:       true,
			},
		},
	}

	attrs := map[string]bool{}

	for _, attrString := range attrStrings {
		split := strings.Split(attrString, ":")
		if len(split) < 2 {
			return nil, fmt.Errorf("invalid attribute format: %s", attrString)
		}
		attrName := split[0]
		attrType := split[1]

		attrSchema, err := common.NewAttributeSchema(attrType, attrType)
		if err != nil {
			return nil, err
		}

		for _, option := range split[2:] {
			switch {
			case option == "notshown":
				attrSchema.ShownInTable = false
			case option == "required":
				attrSchema.Required = true
			case regexp.MustCompile(`^max=([-0-9]+)$`).MatchString(option):
				matches := regexp.MustCompile(`^max=([-0-9]+)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, err
				}
				attrSchema.Max = val
			case regexp.MustCompile(`^min=([-0-9]+)$`).MatchString(option):
				matches := regexp.MustCompile(`^min=([-0-9]+)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, err
				}
				attrSchema.Min = val
			case regexp.MustCompile(`^maxlength=(\d+)$`).MatchString(option):
				matches := regexp.MustCompile(`^maxlength=(\d+)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				attrSchema.MaxLength = val
			case regexp.MustCompile(`^minlength=(\d+)$`).MatchString(option):
				matches := regexp.MustCompile(`^minlength=(\d+)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				attrSchema.MinLength = val
			case option == "private":
				attrSchema.Private = true
			case option == "noteditable":
				attrSchema.Editable = false
			case regexp.MustCompile(`^default=(.+)=$`).MatchString(option):
				matches := regexp.MustCompile(`^default=(.+)$`).FindStringSubmatch(option)
				attrSchema.Default = matches[1]
			case option == "notnullable":
				attrSchema.Nullable = false
			case option == "unique":
				attrSchema.Unique = true
			default:
				return nil, fmt.Errorf("invalid option: %s", option)
			}
		}

		if _, ok := attrs[attrName]; ok {
			return nil, fmt.Errorf("duplicate attribute: %s", attrName)
		}
		attrs[attrName] = true

		modelSchema.Attributes = append(modelSchema.Attributes, attrSchema)

	}

	return modelSchema, nil
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

func ctxFromSchema(schema *common.ModelSchema) (*ModelCtx, error) {
	caser := cases.Title(language.English)

	modfile, err := getModfile()
	if err != nil {
		return nil, err
	}

	modelCtx := &ModelCtx{
		ModFile:           modfile,
		CollectionName:    schema.CollectionName,
		SingularName:      caser.String(schema.SingularName),
		SingularNameCamel: schema.SingularName,
		PluralName:        caser.String(schema.PluralName),
		Attributes:        []*AttributeCtx{},
	}

	for _, attr := range schema.Attributes {
		attrCtx := &AttributeCtx{
			NameCamel:  attr.Name,
			NamePascal: caser.String(attr.Name),
		}

		switch attr.SimpleType {
		case "Number":
			attrCtx.TypeLower = "int"
			attrCtx.TypeUpper = "Int"
		case "String":
			attrCtx.TypeLower = "string"
			attrCtx.TypeUpper = "String"
		case "Boolean":
			attrCtx.TypeLower = "bool"
			attrCtx.TypeUpper = "Bool"
		}

		modelCtx.Attributes = append(modelCtx.Attributes, attrCtx)
	}

	return modelCtx, nil
}
