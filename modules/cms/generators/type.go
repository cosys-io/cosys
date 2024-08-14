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

	if err = generateApi(typesDir, controllersDir, routesDir, ctx); err != nil {
		return err
	}

	return nil
}

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

func generateModel(typeSnakeName string, typesDir string, schema *schema.ModelSchema, ctx *modelCtx) error {
	typeDir := filepath.Join(typesDir, typeSnakeName)
	generator := gen.NewGenerator(
		gen.NewDir(typeDir, gen.GenHeadOnly),
		gen.NewFile(filepath.Join(typeDir, "schema.yaml"), SchemaYamlTmpl, schema),
		gen.NewFile(filepath.Join(typeDir, "schema.go"), SchemaGoTmpl, schema),
		gen.NewFile(filepath.Join(typeDir, "model.go"), ModelTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

func generateApi(typesDir, controllersDir, routesDir string, ctx *modelCtx) error {
	generator := gen.NewGenerator(
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `import \(`, ModelsImportTmpl, ctx),
		gen.ModifyFile(filepath.Join(typesDir, "models.go"), `var Models = map\[string\]common\.Model\{`, ModelsStructTmpl, ctx),
		gen.NewFile(filepath.Join(controllersDir, ctx.PluralSnakeName+"_controllers.go"), ControllerStructTmpl, ctx),
		gen.ModifyFile(filepath.Join(controllersDir, "controllers.go"), `var Controllers = \[\]common\.Controller\{`, ControllersStructTmpl, ctx),
		gen.ModifyFile(filepath.Join(routesDir, "routes.go"), `var Routes = \[\]common\.Route\{`, RoutesStructTmpl, ctx),
	)
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

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

type attrCtx struct {
	TypeLower  string
	TypeUpper  string
	CamelName  string
	PascalName string
}

var ModelTmpl = `package {{.PluralCamelName}}
	
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

var SchemaGoTmpl = `package {{.PluralName}}

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
    editable: false{{end}}{{if .Enum}}
	enum: {{range .Enum}}
		- {{.}}{{end}}{{end}}{{if .Default}}
    default: {{.Default}}{{end}}{{if not .Nullable}}
    nullable: false{{end}}{{if .Unique}}
    unique: true{{end}}
{{end}}`

var ModelsImportTmpl = `import (
	"{{.ModFile}}/{{.TypesDir}}/{{.PluralSnakeName}}"`

var ModelsStructTmpl = `var Models = map[string]common.Model{
	"api.{{.PluralCamelName}}": {{.PluralCamelName}}.{{.PluralPascalName}},`

var ControllersStructTmpl = `var Controllers = []common.Controller{
	{{.PluralCamelName}}Controller,`

var ControllerStructTmpl = `package controllers

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

var RoutesStructTmpl = `var Routes = []common.Route{
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}`" + `, common.GetAction("{{.PluralCamelName}}.findMany")),
	common.NewRoute("GET", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.findOne")),
	common.NewRoute("POST", ` + "`/api/{{.PluralKebabName}}`" + `, common.GetAction("{{.PluralCamelName}}.create")),
	common.NewRoute("PUT", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.update")),
	common.NewRoute("DELETE", ` + "`/api/{{.PluralKebabName}}/{id}`" + `, common.GetAction("{{.PluralCamelName}}.delete")),`
