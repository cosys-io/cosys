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
	Name               string `yaml:"name" json:"name"`
	SimplifiedDataType string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedDataType   string `yaml:"detailedDataType" json:"detailedDataType"`

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
	if attrName == "" {
		return nil, fmt.Errorf("invalid attribute name: %s", attrName)
	}

	var simpleType string
	var detailedType string

	switch attrType {
	case "string":
		simpleType = "String"
		detailedType = "String"
	case "int":
		simpleType = "Number"
		detailedType = "Int"
	case "float":
		simpleType = "Number"
		detailedType = "Float"
	case "boolean":
		simpleType = "Boolean"
		detailedType = "Boolean"
	case "date":
		simpleType = "Date"
		detailedType = "Date"
	case "datetime":
		simpleType = "DateTime"
		detailedType = "DateTime"
	case "timestamp":
		simpleType = "TimeStamp"
		detailedType = "TimeStamp"
	default:
		return nil, fmt.Errorf("invalid attribute type: %s", attrType)
	}

	return &AttributeSchema{
		Name:               attrName,
		SimplifiedDataType: simpleType,
		DetailedDataType:   detailedType,
		ShownInTable:       true,
		Required:           false,
		Max:                2147483647,
		Min:                -2147483648,
		MaxLength:          -1,
		MinLength:          -1,
		Private:            false,
		Editable:           true,
		Default:            "",
		Nullable:           true,
		Unique:             false,
	}, nil
}

var IdSchema = AttributeSchema{
	Name:               "id",
	SimplifiedDataType: "Number",
	DetailedDataType:   "Int",
	ShownInTable:       true,
	Required:           true,
	Max:                2147483647,
	Min:                -2147483648,
	MaxLength:          -1,
	MinLength:          -1,
	Private:            false,
	Editable:           false,
	Nullable:           false,
	Unique:             true,
}

func GetSchema(path string) (*ModelSchema, error) {
	schemaParsed := ModelSchemaParsed{}

	if err := ParseFile(path, &schemaParsed, false); err != nil {
		return nil, err
	}

	return schemaParsed.Schema()
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

func (m ModelSchemaParsed) Schema() (*ModelSchema, error) {
	if m.DisplayName == "" {
		return nil, fmt.Errorf("model has no display name")
	}
	if m.ModelType == "" {
		return nil, fmt.Errorf("model has no model type: %s", m.DisplayName)
	}
	if m.CollectionName == "" {
		return nil, fmt.Errorf("model has no collection name: %s", m.DisplayName)
	}
	if m.SingularName == "" {
		return nil, fmt.Errorf("model has no singular name: %s", m.DisplayName)
	}
	if m.PluralName == "" {
		return nil, fmt.Errorf("model has no plural name: %s", m.DisplayName)
	}

	var attrs []*AttributeSchema

	for _, attr := range m.Attributes {
		attrSchema, err := attr.Schema()
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, attrSchema)
	}

	return &ModelSchema{
		ModelType:      m.ModelType,
		CollectionName: m.CollectionName,
		DisplayName:    m.DisplayName,
		SingularName:   m.SingularName,
		PluralName:     m.PluralName,
		Description:    m.Description,
		Attributes:     attrs,
	}, nil
}

type AttributeSchemaParsed struct {
	Name     string `yaml:"name" json:"name"`
	DataType string `yaml:"dataType" json:"dataType"`

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

func (a AttributeSchemaParsed) Schema() (*AttributeSchema, error) {
	if a.Name == "" {
		return nil, fmt.Errorf("attribute has no name")
	}
	if a.DataType == "" {
		return nil, fmt.Errorf("attribute has no type: %s", a.Name)
	}

	var simplifiedDataType string
	var detailedDataType string
	switch a.DataType {
	case "Int":
		simplifiedDataType = "Number"
		detailedDataType = "Int"
	case "Float":
		simplifiedDataType = "Number"
		detailedDataType = "Float"
	case "String":
		simplifiedDataType = "String"
		detailedDataType = "String"
	case "Boolean":
		simplifiedDataType = "Boolean"
		detailedDataType = "Boolean"
	case "Date":
		simplifiedDataType = "Date"
		detailedDataType = "Date"
	case "DateTime":
		simplifiedDataType = "DateTime"
		detailedDataType = "DateTime"
	case "TimeStamp":
		simplifiedDataType = "TimeStamp"
		detailedDataType = "TimeStamp"
	default:
		return nil, fmt.Errorf("invalid attribute type: %s", a.DataType)
	}

	return &AttributeSchema{
		Name:               a.Name,
		SimplifiedDataType: simplifiedDataType,
		DetailedDataType:   detailedDataType,

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
		Unique:   checkDefault(false, a.Unique),
	}, nil
}

func checkDefault[T any](defaultValue T, value *T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
