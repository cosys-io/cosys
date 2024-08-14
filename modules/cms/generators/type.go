package generators

import (
	"bufio"
	"github.com/cosys-io/cosys/common"
	gen "github.com/cosys-io/cosys/cosys_cli/generator"
	"github.com/cosys-io/cosys/modules/cms/schema"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// GenerateType generates the files for a new collection type.
func GenerateType(schema *schema.ModelSchema) error {
	common.InitConfigs()

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

	ctx, err := getCtx(schema, typesDir)
	if err != nil {
		return err
	}

	typeSnakeName := strcase.ToSnake(schema.PluralName())

	if err = generateModel(typeSnakeName, typesDir, schema, ctx); err != nil {
		return err
	}

	if err = generateApi(controllersDir, routesDir, ctx); err != nil {
		return err
	}

	return nil
}

// getCtx returns the code generator context from the given model schema.
func getCtx(modelSchema *schema.ModelSchema, typesDir string) (*modelCtx, error) {
	modFile, err := getModFile()
	if err != nil {
		return nil, err
	}

	ctx := getModelCtx(modelSchema, modFile, typesDir)

	for index, attrSchema := range modelSchema.Attributes() {
		ctx.Attributes[index] = getAttrCtx(attrSchema.(*schema.AttributeSchema))
	}

	return ctx, nil
}

// getModFile returns the module name.
func getModFile() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		return line[7:], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("module not found")
}

// getModelCtx returns the code generator context from the given model schema, without the attribute contexts.
func getModelCtx(schema *schema.ModelSchema, modFile string, typesDir string) *modelCtx {
	return &modelCtx{
		DBName:             schema.CollectionName(),
		DisplayName:        schema.DisplayName(),
		SingularCamelName:  schema.SingularName(),
		PluralCamelName:    schema.PluralName(),
		SingularPascalName: strcase.ToCamel(schema.SingularName()),
		PluralPascalName:   strcase.ToCamel(schema.PluralName()),
		SingularSnakeName:  strcase.ToSnake(schema.SingularName()),
		PluralSnakeName:    strcase.ToSnake(schema.PluralName()),
		SingularKebabName:  strcase.ToKebab(schema.SingularName()),
		PluralKebabName:    strcase.ToKebab(schema.PluralName()),
		SingularHumanName:  strcase.ToDelimited(schema.SingularName(), ' '),
		PluralHumanName:    strcase.ToDelimited(schema.PluralName(), ' '),

		ModFile:    modFile,
		TypesDir:   typesDir,
		Attributes: make([]*attrCtx, len(schema.Attributes())),
	}
}

// getAttrCtx returns the code generator context from the given attribute schema.
func getAttrCtx(schema *schema.AttributeSchema) *attrCtx {
	ctx := &attrCtx{
		CamelName:  schema.Name(),
		PascalName: strcase.ToCamel(schema.Name()),
	}

	switch schema.SimplifiedDataType() {
	case "Number":
		ctx.TypeLower = "int"
		ctx.TypeUpper = "Int"
	case "String":
		ctx.TypeLower = "string"
		ctx.TypeUpper = "String"
	case "Boolean":
		ctx.TypeLower = "bool"
		ctx.TypeUpper = "Bool"
	}

	return ctx
}

// generateModel generates the code for the model of the collection type.
func generateModel(typeSnakeName string, typesDir string, schema *schema.ModelSchema, ctx *modelCtx) error {
	typeDir := filepath.Join(typesDir, typeSnakeName)
	generator := gen.NewGenerator(
		gen.NewDir(typeDir, gen.GenHeadOnly),
		gen.NewFile(filepath.Join(typeDir, "schema.yaml"), schemaYamlTmpl, schema),
		gen.NewFile(filepath.Join(typeDir, "schema.go"), schemaGoTmpl, schema),
		gen.NewFile(filepath.Join(typeDir, "model.go"), modelTmpl, ctx),
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `import \(`, modelsImportTmpl, ctx),
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `var Models = map\[string\]common\.Model\{`, modelsStructTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// generateApi generates the code for the controllers and routes of the collection type.
func generateApi(controllersDir, routesDir string, ctx *modelCtx) error {
	generator := gen.NewGenerator(
		gen.NewFile(filepath.Join(controllersDir, ctx.PluralSnakeName+"_controllers.go"), typeControllerTmpl, ctx),
		gen.ModifyFile(filepath.Join(controllersDir, "controllers.go"), `var Controllers = \[\]common\.Controller\{`, controllersTmpl, ctx),
		gen.ModifyFile(filepath.Join(routesDir, "routes.go"), `var Routes = \[\]common\.Route\{`, routesTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// modelCtx contains the data for generating code for a new collection type.
type modelCtx struct {
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
	Attributes []*attrCtx
}

// attrCtx contains the data for generating code for an attribute of a new collection type.
type attrCtx struct {
	TypeLower  string
	TypeUpper  string
	CamelName  string
	PascalName string
}

// modelTmpl is the template for creating the entity type, model type and model for a new collection type.
var modelTmpl = `package {{.PluralCamelName}}
	
import (
	"github.com/cosys-io/cosys/common"
)

type {{.SingularPascalName}} struct {
{{range .Attributes}}    {{.PascalName}} {{.TypeLower}} ` + "`" + `json:"{{.CamelName}}"` + "`" + `
{{end}}}


type {{.PluralPascalName}}Model struct {
	*common.ModelBase

{{range .Attributes}}    {{.PascalName}} common.{{.TypeUpper}}Attribute
{{end}}}

var {{.PluralPascalName}}, _ = common.NewModel[{{.SingularPascalName}}, {{.PluralPascalName}}Model](
	"{{.DBName}}", 
	"{{.SingularCamelName}}", 
	"{{.PluralCamelName}}", 
	{{.SingularCamelName}}Schema,
)`

// schemaGoTmpl is the template for creating the schema struct of a new collection type.
var schemaGoTmpl = `package {{.PluralName}}

import (
	"github.com/cosys-io/cosys/modules/cms/schema"
)

var {{.SingularName}}Schema = schema.NewModelSchema(
	"{{.CollectionName}}",
	"{{.DisplayName}}",
	"{{.SingularName}}",
	"{{.PluralName}}",
	"{{.Description}}",
{{range .Attributes}}        schema.NewAttrSchema(
			"{{.Name}}",
			"{{.SimplifiedDataType}}",
			"{{.DetailedDataType}}",{{if not .ShownInTable}}
			schema.HideInTable,{{end}}{{if .Required}}
			schema.Required,{{end}}{{if ne .Max 2147483647}}
			schema.Max({{.Max}}),{{end}}{{if ne .Min -2147483648}}
			schema.Min({{.Min}}),{{end}}{{if ne .MaxLength -1}}
			schema.MaxLength({{.MaxLength}}),{{end}}{{if ne .MinLength -1}}
			schema.MinLength({{.MinLength}}),{{end}}{{if .Private}}
			schema.Private,{{end}}{{if not .Editable}}
			schema.NotEditable,{{end}}{{if .Enum}}
			schema.Enum([]string{ {{range.Enum}}
				"{{.}}",{{end}}
			},{{end}}{{if .Default}}
			schema.Default("{{.Default}}"),{{end}}{{if not .Nullable}}
			schema.NotNullable,{{end}}{{if .Unique}}
			schema.Unique,{{end}}
		),
{{end}})
`

// schemaYamlTmpl is the template for creating the yaml configuration of a new collection type.
var schemaYamlTmpl = `modelType: {{.ModelType}}
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
    editable: false{{end}}{{if .Enum}}
	enum: {{range .Enum}}
		- {{.}}{{end}}{{end}}{{if .Default}}
    default: {{.Default}}{{end}}{{if not .Nullable}}
    nullable: false{{end}}{{if .Unique}}
    unique: true{{end}}
{{end}}`

// modelsImportTmpl is the template for adding the import for the model of a new collection type
// to the imports in the models.go file.
var modelsImportTmpl = `import (
	"{{.ModFile}}/{{.TypesDir}}/{{.PluralSnakeName}}"`

// modelsStructTmpl is the template for the adding the model of a new collection type
// to the models map in the models.go file.
var modelsStructTmpl = `var Models = map[string]common.Model{
	"api.{{.PluralCamelName}}": {{.PluralCamelName}}.{{.PluralPascalName}},`

// controllersTmpl is the template for adding the controller of a new collection type
// to the controllers slice in the controllers.go file.
var controllersTmpl = `var Controllers = []common.Controller{
	{{.PluralCamelName}}Controller,`

// typeControllerTmpl is the template for creating the controller of a new collection type
// in a new file in the controllers package.
var typeControllerTmpl = `package controllers

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/cms/routes"
)

var {{.PluralCamelName}}Controller, _ = common.NewController("{{.PluralCamelName}}", map[string]common.ActionFunc{
	"findMany": routes.FindMany("api.{{.PluralCamelName}}"),
	"findOne": routes.FindOne("api.{{.PluralCamelName}}"),
	"create": routes.Create("api.{{.PluralCamelName}}"),
	"update": routes.Update("api.{{.PluralCamelName}}"),
	"delete": routes.Delete("api.{{.PluralCamelName}}"),
})`

// routesTmpl is the template for adding the routes for a new collection type
// to the routes slice in the routes.go file.
var routesTmpl = `var Routes = []common.Route{
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}`" + `, common.GetAction("{{.PluralCamelName}}.findMany")),
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.findOne")),
	common.NewRoute("POST", ` + "`/api/{{.PluralKebabName}}`" + `, common.GetAction("{{.PluralCamelName}}.create")),
	common.NewRoute("PUT", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.update")),
	common.NewRoute("DELETE", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.delete")),`
