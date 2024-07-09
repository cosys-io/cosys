package cms

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	gen "github.com/cosys-io/cosys/cosys_cli/cmd/generator"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
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
		common.InitConfigs()

		collectionName := args[0]
		attributes := args[1:]

		schema, err := schemaFromArgs(collectionName, displayName, singularName, pluralName, description, attributes)
		if err != nil {
			log.Fatal(err)
		}

		if err := GenerateType(schema); err != nil {
			log.Fatal(err)
		}
	},
}

func GenerateType(schema *common.ModelSchema) error {
	typesDir, err := common.GetPathConfig("cms_content_types_path", true)
	if err != nil {
		return err
	}
	routesDir, err := common.GetPathConfig("cms_routes_path", true)
	if err != nil {
		return err
	}
	controllersDir, err := common.GetPathConfig("cms_controllers_path", true)
	if err != nil {
		return err
	}

	ctx, err := ctxFromSchema(schema, typesDir)
	if err != nil {
		return err
	}

	typeSnakeName := strcase.ToSnake(schema.PluralName)

	if err = generateModel(typeSnakeName, typesDir, schema, ctx); err != nil {
		return err
	}

	if err = generateApi(typeSnakeName, typesDir, routesDir, controllersDir, ctx); err != nil {
		return err
	}

	return nil
}

func generateModel(typeSnakeName string, typesDir string, schema *common.ModelSchema, ctx *ModelCtx) error {
	typeDir := filepath.Join(typesDir, typeSnakeName)
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

func generateApi(typeSnakeName, typesDir, routesDir, controllersDir string, ctx *ModelCtx) error {
	generator := gen.NewGenerator(
		gen.NewFile(filepath.Join(controllersDir, typeSnakeName+"_controllers.go"), ModelControllerTmpl, ctx),
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `import \(`, ModelsImportTmpl, ctx),
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `var Models = map\[string\]common\.Model\{`, ModelsStructTmpl, ctx),
		gen.ModifyFile(filepath.Join(controllersDir, "controllers.go"), `var Controllers = map\[string\]common\.Controller\{`, ControllersStructTmpl, ctx),
		gen.ModifyFile(filepath.Join(routesDir, "routes.go"), `var Routes = \[\]\*common\.Route\{`, RoutesStructTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

func schemaFromArgs(collectionName string, displayName string, singularName string, pluralName string,
	description string, attrStrings []string) (*common.ModelSchema, error) {
	modelSchema := &common.ModelSchema{
		ModelType:      "collectionType",
		CollectionName: strcase.ToLowerCamel(collectionName),
		DisplayName:    displayName,
		SingularName:   strcase.ToLowerCamel(singularName),
		PluralName:     strcase.ToLowerCamel(pluralName),
		Description:    description,
		Attributes: []*common.AttributeSchema{
			&common.IdSchema,
		},
	}

	attrs := make(map[string]bool)

	for _, attrString := range attrStrings {
		split := strings.Split(attrString, ":")
		if len(split) < 2 {
			return nil, fmt.Errorf("invalid attribute format: %s", attrString)
		}
		attrName := strcase.ToLowerCamel(split[0])
		attrType := split[1]

		attrSchema, err := common.NewAttributeSchema(attrName, attrType)
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

		if _, dup := attrs[attrName]; dup {
			return nil, fmt.Errorf("duplicate attribute: %s", attrName)
		}
		attrs[attrName] = true

		modelSchema.Attributes = append(modelSchema.Attributes, attrSchema)

	}

	return modelSchema, nil
}

type ModelCtx struct {
	DBName             string
	DisplayName        string
	SingularCamelName  string
	PluralCamelName    string
	SingularPascalName string
	PluralPascalName   string
	SingularSnakeName  string
	PluralSnakeName    string
	SingularKebabName  string
	PluralKebabName    string
	SingularHumanName  string
	PluralHumanName    string

	ModFile    string
	TypesDir   string
	Attributes []*AttributeCtx
}

type AttributeCtx struct {
	TypeLower  string
	TypeUpper  string
	CamelName  string
	PascalName string
}

func ctxFromSchema(schema *common.ModelSchema, typesDir string) (*ModelCtx, error) {
	modFile, err := getModFile()
	if err != nil {
		return nil, err
	}

	modelCtx := ModelCtx{
		DBName:             schema.CollectionName,
		DisplayName:        schema.DisplayName,
		SingularCamelName:  schema.SingularName,
		PluralCamelName:    schema.PluralName,
		SingularPascalName: strcase.ToCamel(schema.SingularName),
		PluralPascalName:   strcase.ToCamel(schema.PluralName),
		SingularSnakeName:  strcase.ToSnake(schema.SingularName),
		PluralSnakeName:    strcase.ToSnake(schema.PluralName),
		SingularKebabName:  strcase.ToKebab(schema.SingularName),
		PluralKebabName:    strcase.ToKebab(schema.PluralName),
		SingularHumanName:  strcase.ToDelimited(schema.SingularName, ' '),
		PluralHumanName:    strcase.ToDelimited(schema.PluralName, ' '),

		ModFile:    modFile,
		TypesDir:   typesDir,
		Attributes: []*AttributeCtx{},
	}

	for _, attr := range schema.Attributes {
		attrCtx := AttributeCtx{
			CamelName:  attr.Name,
			PascalName: strcase.ToCamel(attr.Name),
		}

		switch attr.SimplifiedDataType {
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

		modelCtx.Attributes = append(modelCtx.Attributes, &attrCtx)
	}

	return &modelCtx, nil
}

var ModelTmpl = `package {{.PluralCamelName}}
	
import (
	"log"
	"github.com/cosys-io/cosys/common"
)

var (
	Schema = &common.ModelSchema{}
)

func init() {
	var err error
	Schema, err = common.GetSchema("{{.TypesDir}}/{{.PluralSnakeName}}/schema.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

type {{.SingularPascalName}} struct {
{{range .Attributes}}    {{.PascalName}} {{.TypeLower}} ` + "`" + `json:"{{.CamelName}}"` + "`" + `
{{end}}}


type {{.PluralPascalName}}Model struct {
	common.ModelBase
	lifecycle common.Lifecycle

{{range .Attributes}}    {{.PascalName}} *common.{{.TypeUpper}}Attribute
{{end}}}

var {{.PluralPascalName}} = {{.PluralPascalName}}Model{
	common.NewModelBase("{{.DBName}}", "{{.DisplayName}}", "{{.SingularCamelName}}", "{{.PluralCamelName}}"),
	Lifecycle,

{{range .Attributes}}    common.New{{.TypeUpper}}Attribute("{{.CamelName}}"),
{{end}}}

func (m {{.PluralPascalName}}Model) New_() common.Entity {
	return &{{.SingularPascalName}}{}
}

{{$Model := .PluralPascalName}}
func (m {{.PluralPascalName}}Model) All_() []common.Attribute {
	return []common.Attribute{
{{range .Attributes}}        {{$Model}}.{{.PascalName}},
{{end}}}
}

func (m {{.PluralPascalName}}Model) Id_() *common.IntAttribute {
	return {{.PluralPascalName}}.Id
}

func (m {{.PluralPascalName}}Model) Schema_() *common.ModelSchema {
	return Schema
}

func (m {{.PluralPascalName}}Model) Lifecycle_() common.Lifecycle {
	return m.lifecycle
}`

var LifecycleTmpl = `package {{.PluralCamelName}}

import "github.com/cosys-io/cosys/common"

var Lifecycle = common.NewLifecycle()`

var SchemaYamlTmpl = `modelType: {{.ModelType}}
collectionName: {{.CollectionName}}
displayName: {{.DisplayName}}
singularName: {{.SingularName}}
pluralName: {{.PluralName}}
description: {{.Description}}
attributes:
{{range .Attributes}}  - name: {{.Name}}
    simplifiedDataType: {{.SimplifiedDataType}}
    detailedDataType: {{.DetailedDataType}}{{if not .ShownInTable}}
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
	"strings"

	"github.com/cosys-io/cosys/common"
)

var {{.PluralPascalName}}Controller = map[string]common.Action{
	"findMany": findMany{{.PluralPascalName}},
	"findOne": findOne{{.SingularPascalName}},
	"create":  create{{.SingularPascalName}},
	"update":  update{{.SingularPascalName}},
	"delete":  delete{{.SingularPascalName}},
}

func findMany{{.PluralPascalName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		model, ok := cs.Models["api.{{.PluralCamelName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		attrSlice := model.All_()

		page := 1
		pageSize := int64(20)
		sort := []*common.Order{}
		filter := []common.Condition{}
		fields := []common.Attribute{}
		populate := []common.Attribute{}

		pageSizeString, ok := params["pageSize"]
		if ok {
			pageSize, err = strconv.ParseInt(pageSizeString, 10, 64)
			if err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}
		}

		pageString, ok := params["page"]
		if ok {
			page, err = strconv.Atoi(pageString)
			if err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}
		}

		sortSliceString, ok := params["sort"]
		if ok {
			sortSlice := strings.Split(sortSliceString, ",")
			for _, sortString := range sortSlice {
				if len(sortString) == 0 {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}

				isAsc := true
				if sortString[0] == '-' {
					isAsc = false
					sortString = sortString[1:]
				}

				var sortAttr common.Attribute

				for _, attr := range attrSlice {
					if attr.CamelName() == sortString {
						sortAttr = attr
					}
				}

				if sortAttr == nil {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}

				if isAsc {
					sort = append(sort, sortAttr.Asc())
				} else {
					sort = append(sort, sortAttr.Desc())
				}
			}
		}

		fieldSliceString, ok := params["fields"]
		if ok {
			fieldSlice := strings.Split(fieldSliceString, ",")
			for _, fieldString := range fieldSlice {
				var fieldAttr common.Attribute

				for _, attr := range attrSlice {
					if attr.CamelName() == fieldString {
						fieldAttr = attr
					}
				}

				if fieldAttr == nil {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}

				fields = append(fields, fieldAttr)
			}
		}

		populateSliceString, ok := params["populate"]
		if ok {
			populateSlice := strings.Split(populateSliceString, ",")
			for _, populateString := range populateSlice {
				var populateAttr common.Attribute

				for _, attr := range attrSlice {
					if attr.CamelName() == populateString {
						populateAttr = attr
					}
				}

				if populateAttr == nil {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}

				populate = append(populate, populateAttr)
			}
		}

		dbParams := common.NewDBParamsBuilder().
			Offset(pageSize * (int64(page) - 1)).
			Limit(pageSize).
			OrderBy(sort...).
			Where(filter...).
			Select(fields...).
			Populate(populate...).
			Build()
		entities, err := cs.Database().FindMany("api.{{.PluralCamelName}}", dbParams)
		if err != nil {
			common.RespondError(w, "Could not find {{.PluralHumanName}}.", http.StatusBadRequest)
			return
		}

		common.RespondMany(w, entities, page, http.StatusOK)
	}
}

func findOne{{.SingularPascalName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		idString, ok := params["documentId"]
		if !ok {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idString)
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		model, ok := cs.Models["api.{{.PluralCamelName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		
		dbParams := common.NewDBParamsBuilder().
			Where(model.Id_().Eq(id)).
			Build()

		entity, err := cs.Database().FindOne("api.{{.PluralCamelName}}", dbParams)
		if err != nil {
			common.RespondError(w, "Could not find {{.SingularHumanName}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, entity, http.StatusOK)
	}
}

func create{{.SingularPascalName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model, ok := cs.Models["api.{{.PluralCamelName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		entity := model.New_()

		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		newEntity, err := cs.Database().Create("api.{{.PluralCamelName}}", entity, common.NewDBParams())
		if err != nil {
			common.RespondError(w, "Could not create {{.SingularHumanName}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, newEntity, http.StatusOK)
	}
}

func update{{.SingularPascalName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		idString, ok := params["documentId"]
		if !ok {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idString)
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		model, ok := cs.Models["api.{{.PluralCamelName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}
		entity := model.New_()

		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		dbParams := common.NewDBParamsBuilder().
			Where(model.Id_().Eq(id)).
			Build()

		newEntity, err := cs.Database().Update("api.{{.PluralCamelName}}", entity, dbParams)
		if err != nil {
			common.RespondError(w, "Could not update {{.SingularHumanName}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, newEntity, http.StatusOK)
	}
}

func delete{{.SingularPascalName}}(cs common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
			return
		}

		idString, ok := params["documentId"]
		if !ok {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idString)
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		model, ok := cs.Models["api.{{.PluralCamelName}}"]
		if !ok {
			common.RespondInternalError(w)
			return
		}

		dbParams := common.NewDBParamsBuilder().
			Where(model.Id_().Eq(id)).
			Build()

		oldEntity, err := cs.Database().Delete("api.{{.PluralCamelName}}", dbParams)
		if err != nil {
			common.RespondError(w, "Could not delete {{.SingularHumanName}}.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, oldEntity, http.StatusOK)
	}
}`

var ModelsImportTmpl = `import (
	"{{.ModFile}}/{{.TypesDir}}/{{.PluralSnakeName}}"`

var ModelsStructTmpl = `var Models = map[string]common.Model{
	"api.{{.PluralCamelName}}": {{.PluralCamelName}}.{{.PluralPascalName}},`

var ControllersStructTmpl = `var Controllers = map[string]common.Controller{
	"{{.PluralCamelName}}": {{.PluralPascalName}}Controller,`

var RoutesStructTmpl = `var Routes = []*common.Route{
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}`" + `, "{{.PluralCamelName}}.findMany"),
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}/{documentId}`" + `, "{{.PluralCamelName}}.findOne"),
	common.NewRoute("POST", ` + "`/api/{{.PluralKebabName}}`" + `, "{{.PluralCamelName}}.create"),
	common.NewRoute("PUT", ` + "`/api/{{.PluralKebabName}}/{documentId}`" + `, "{{.PluralCamelName}}.update"),
	common.NewRoute("DELETE", ` + "`/api/{{.PluralKebabName}}/{documentId}`" + `, "{{.PluralCamelName}}.delete"),`
