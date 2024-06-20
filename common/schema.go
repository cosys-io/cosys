package common

import "fmt"

type ModelSchema struct {
	ModelType      string             `yaml:"modelType" json:"modelType"`
	CollectionName string             `yaml:"collectionName" json:"collectionName"`
	DisplayName    string             `yaml:"displayName" json:"displayName"`
	SingularName   string             `yaml:"singularName" json:"singularName"`
	PluralName     string             `yaml:"pluralName" json:"pluralName"`
	Description    string             `yaml:"description" json:"description"`
	Attributes     []*AttributeSchema `yaml:"attributes" json:"attributes"`
}

type AttributeSchema struct {
	Name         string `yaml:"name" json:"name"`
	SimpleType   string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedType string `yaml:"detailedDataType" json:"detailedDataType"`

	ShownInTable bool  `yaml:"shownInTable" json:"shownInTable"`
	Required     bool  `yaml:"required" json:"required"`
	Max          int64 `yaml:"max" json:"max"`
	Min          int64 `yaml:"min" json:"min"`
	MaxLength    int   `yaml:"maxLength" json:"maxLength"`
	MinLength    int   `yaml:"minLength" json:"minLength"`
	Private      bool  `yaml:"private" json:"private"`
	Editable     bool  `yaml:"editable" json:"editable"`

	Default  string `yaml:"default" json:"default"`
	Nullable bool   `yaml:"nullable" json:"nullable"`
	Unique   bool   `yaml:"unique" json:"unique"`
}

func NewAttributeSchema(attrName, attrType string) (*AttributeSchema, error) {
	attrSchema := &AttributeSchema{
		Name:         attrName,
		SimpleType:   "",
		DetailedType: "",
		ShownInTable: true,
		Required:     false,
		Max:          2147483647,
		Min:          -2147483648,
		MaxLength:    -1,
		MinLength:    -1,
		Private:      false,
		Editable:     true,
		Default:      "",
		Nullable:     true,
		Unique:       false,
	}

	switch attrType {
	case "string":
		attrSchema.SimpleType = "String"
		attrSchema.DetailedType = "String"
	case "int":
		attrSchema.SimpleType = "Number"
		attrSchema.DetailedType = "Int"
	case "float":
		attrSchema.SimpleType = "Number"
		attrSchema.DetailedType = "Float"
	case "boolean":
		attrSchema.SimpleType = "Boolean"
		attrSchema.DetailedType = "Boolean"
	case "date":
		attrSchema.SimpleType = "Date"
		attrSchema.DetailedType = "Date"
	case "datetime":
		attrSchema.SimpleType = "DateTime"
		attrSchema.DetailedType = "DateTime"
	case "timestamp":
		attrSchema.SimpleType = "TimeStamp"
		attrSchema.DetailedType = "TimeStamp"
	default:
		return nil, fmt.Errorf("invalid attribute type: %s", attrType)
	}

	return attrSchema, nil
}

func GetSchema(path string) (*ModelSchema, error) {
	schemaParsed := &ModelSchemaParsed{}

	if err := ParseFile(path, schemaParsed, false); err != nil {
		return nil, err
	}

	return schemaParsed.Schema(), nil
}

type ModelSchemaParsed struct {
	ModelType      string                   `yaml:"modelType" json:"modelType"`
	CollectionName string                   `yaml:"collectionName" json:"collectionName"`
	DisplayName    string                   `yaml:"displayName" json:"displayName"`
	SingularName   string                   `yaml:"singularName" json:"singularName"`
	PluralName     string                   `yaml:"pluralName" json:"pluralName"`
	Description    string                   `yaml:"description" json:"description"`
	Attributes     []*AttributeSchemaParsed `yaml:"attributes" json:"attributes"`
}

func (m ModelSchemaParsed) Schema() *ModelSchema {
	attrs := []*AttributeSchema{}

	for _, attr := range m.Attributes {
		attrs = append(attrs, attr.Schema())
	}

	return &ModelSchema{
		ModelType:      m.ModelType,
		CollectionName: m.CollectionName,
		DisplayName:    m.DisplayName,
		SingularName:   m.SingularName,
		PluralName:     m.PluralName,
		Description:    m.Description,
		Attributes:     attrs,
	}
}

type AttributeSchemaParsed struct {
	Name         string `yaml:"name" json:"name"`
	SimpleType   string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedType string `yaml:"detailedDataType" json:"detailedDataType"`

	ShownInTable *bool  `yaml:"shownInTable" json:"shownInTable"`
	Required     *bool  `yaml:"required" json:"required"`
	Max          *int64 `yaml:"max" json:"max"`
	Min          *int64 `yaml:"min" json:"min"`
	MaxLength    *int   `yaml:"maxLength" json:"maxLength"`
	MinLength    *int   `yaml:"minLength" json:"minLength"`
	Private      *bool  `yaml:"private" json:"private"`
	Editable     *bool  `yaml:"editable" json:"editable"`

	Default  *string `yaml:"default" json:"default"`
	Nullable *bool   `yaml:"nullable" json:"nullable"`
	Unique   *bool   `yaml:"unique" json:"unique"`
}

func (a AttributeSchemaParsed) Schema() *AttributeSchema {
	return &AttributeSchema{
		Name:         a.Name,
		SimpleType:   a.SimpleType,
		DetailedType: a.DetailedType,

		ShownInTable: checkDefault(true, a.ShownInTable),
		Required:     checkDefault(false, a.Required),
		Max:          checkDefault(int64(2147483647), a.Max),
		Min:          checkDefault(int64(-2147483648), a.Min),
		MaxLength:    checkDefault(-1, a.MaxLength),
		MinLength:    checkDefault(-1, a.MinLength),
		Private:      checkDefault(false, a.Private),
		Editable:     checkDefault(true, a.Editable),

		Default:  checkDefault("", a.Default),
		Nullable: checkDefault(true, a.Nullable),
		Unique:   checkDefault(true, a.Unique),
	}
}

func checkDefault[T any](defaultValue T, value *T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
