package content_builder

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/cosys_cli/cmd"
	"gopkg.in/yaml.v3"
	"net/http"
)

var Controller = map[string]common.Action{
	"schema": schema,
	"build":  build,
}

func schema(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schemas := []common.ModelSchema{}

		for _, model := range cosys.Models {
			schemas = append(schemas, *model.Schema_())
		}

		common.RespondMany(w, schemas, 1, http.StatusOK)
	}
}

func build(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
		}

		name := params[0]

		schema := &ModelSchemaRequest{}

		if err := yaml.NewDecoder(r.Body).Decode(schema); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		if err := cmd.GenerateType(name, schema.Schema()); err != nil {
			common.RespondError(w, "Unable to build content type.", http.StatusBadRequest)
		}

		common.RespondOne(w, "Content type successfully created.", 200)
	}
}

type ModelSchemaRequest struct {
	ModelType      string                    `yaml:"modelType" json:"modelType"`
	CollectionName string                    `yaml:"collectionName" json:"collectionName"`
	DisplayName    string                    `yaml:"displayName" json:"displayName"`
	SingularName   string                    `yaml:"singularName" json:"singularName"`
	PluralName     string                    `yaml:"pluralName" json:"pluralName"`
	Description    string                    `yaml:"description" json:"description"`
	Attributes     []*AttributeSchemaRequest `yaml:"attributes" json:"attributes"`
}

func (m ModelSchemaRequest) Schema() *common.ModelSchema {
	attrs := []*common.AttributeSchema{}

	for _, attr := range m.Attributes {
		attrs = append(attrs, attr.Schema())
	}

	return &common.ModelSchema{
		ModelType:      m.ModelType,
		CollectionName: m.CollectionName,
		DisplayName:    m.DisplayName,
		SingularName:   m.SingularName,
		PluralName:     m.PluralName,
		Description:    m.Description,
		Attributes:     attrs,
	}
}

type AttributeSchemaRequest struct {
	Name         string `yaml:"name" json:"name"`
	SimpleType   string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedType string `yaml:"detailedDataType" json:"detailedDataType"`

	ShownInTable bool   `yaml:"shownInTable" json:"shownInTable"`
	Required     bool   `yaml:"required" json:"required"`
	Max          *int64 `yaml:"max" json:"max"`
	Min          *int64 `yaml:"min" json:"min"`
	MaxLength    *int   `yaml:"maxLength" json:"maxLength"`
	MinLength    *int   `yaml:"minLength" json:"minLength"`
	Private      bool   `yaml:"private" json:"private"`
	Editable     bool   `yaml:"editable" json:"editable"`

	Default  string `yaml:"default" json:"default"`
	Nullable bool   `yaml:"nullable" json:"nullable"`
	Unique   bool   `yaml:"unique" json:"unique"`
}

func (a AttributeSchemaRequest) Schema() *common.AttributeSchema {
	maxVal := int64(2147483647)
	minVal := int64(-2147483648)
	maxLength := -1
	minLength := -1
	if a.Max != nil {
		maxVal = *a.Max
	}
	if a.Min != nil {
		minVal = *a.Min
	}
	if a.MaxLength != nil {
		maxLength = *a.MaxLength
	}
	if a.MinLength != nil {
		minLength = *a.MinLength
	}

	return &common.AttributeSchema{
		Name:         a.Name,
		SimpleType:   a.SimpleType,
		DetailedType: a.DetailedType,

		ShownInTable: a.ShownInTable,
		Required:     a.Required,
		Max:          maxVal,
		Min:          minVal,
		MaxLength:    maxLength,
		MinLength:    minLength,
		Private:      a.Private,
		Editable:     a.Editable,

		Default:  a.Default,
		Nullable: a.Nullable,
		Unique:   a.Unique,
	}
}
