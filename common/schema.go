package common

type ModelSchema struct {
	CollectionName string             `yaml:"collectionName" json:"collectionName"`
	SingularName   string             `yaml:"singularName" json:"singularName"`
	PluralName     string             `yaml:"pluralName" json:"pluralName"`
	Attributes     []*AttributeSchema `yaml:"attributes" json:"attributes"`
}

type AttributeSchema struct {
	Name               string `yaml:"name" json:"name"`
	SimplifiedDataType string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedDataType   string `yaml:"detailedDataType" json:"detailedDataType"`

	Required  bool  `yaml:"required" json:"required"`
	Max       int64 `yaml:"max" json:"max"`
	Min       int64 `yaml:"min" json:"min"`
	MaxLength int   `yaml:"maxLength" json:"maxLength"`
	MinLength int   `yaml:"minLength" json:"minLength"`
	Private   bool  `yaml:"private" json:"private"`
	Editable  bool  `yaml:"editable" json:"editable"`

	Default  string `yaml:"default" json:"default"`
	Nullable bool   `yaml:"nullable" json:"nullable"`
	Unique   bool   `yaml:"unique" json:"unique"`
}

func NewModelSchema(singular, plural string, attrs ...*AttributeSchema) *ModelSchema {
	return &ModelSchema{
		SingularName: singular,
		PluralName:   plural,
		Attributes:   attrs,
	}
}

func NewAttrSchema(attrName, simpleType, detailedType string, opts ...AttrOption) *AttributeSchema {
	schema := &AttributeSchema{
		Name:               attrName,
		SimplifiedDataType: simpleType,
		DetailedDataType:   detailedType,
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
	}

	for _, opt := range opts {
		opt(schema)
	}

	return schema
}

type AttrOption func(*AttributeSchema)

var Required AttrOption = func(schema *AttributeSchema) {
	schema.Required = true
}

func Max(max int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.Max = max
	}
}

func Min(min int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.Min = min
	}
}

func MaxLength(max int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.MaxLength = max
	}
}

func MinLength(min int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.MinLength = min
	}
}

var Private AttrOption = func(schema *AttributeSchema) {
	schema.Private = true
}

var NotEditable AttrOption = func(schema *AttributeSchema) {
	schema.Editable = false
}

func Default(def string) AttrOption {
	return func(schema *AttributeSchema) {
		schema.Default = def
	}
}

var NotNullable AttrOption = func(schema *AttributeSchema) {
	schema.Nullable = false
}

var Unique AttrOption = func(schema *AttributeSchema) {
	schema.Unique = true
}

var IdSchema = AttributeSchema{
	Name:               "id",
	SimplifiedDataType: "Number",
	DetailedDataType:   "Int",
	Required:           false,
	Max:                2147483647,
	Min:                -2147483648,
	MaxLength:          -1,
	MinLength:          -1,
	Private:            false,
	Editable:           false,
	Nullable:           false,
	Unique:             true,
}

var UuidSchema = AttributeSchema{
	Name:               "uuid",
	SimplifiedDataType: "String",
	DetailedDataType:   "String",
	Required:           false,
	Max:                2147483647,
	Min:                -2147483648,
	MaxLength:          -1,
	MinLength:          -1,
	Private:            false,
	Editable:           false,
	Nullable:           false,
	Unique:             true,
}

//func GetSchema(path string) (*ModelSchema, error) {
//	schemaParsed := ModelSchemaParsed{}
//
//	if err := ParseFile(path, &schemaParsed, false); err != nil {
//		return nil, err
//	}
//
//	return schemaParsed.Schema()
//}
//
//type ModelSchemaParsed struct {
//	ModelType      string                   `yaml:"modelType" json:"modelType"`
//	CollectionName string                   `yaml:"collectionName" json:"collectionName"`
//	DisplayName    string                   `yaml:"displayName" json:"displayName"`
//	SingularName   string                   `yaml:"singularName" json:"singularName"`
//	PluralName     string                   `yaml:"pluralName" json:"pluralName"`
//	Description    string                   `yaml:"description" json:"description"`
//	Attributes     []*AttributeSchemaParsed `yaml:"attributes" json:"attributes"`
//}
//
//func (m ModelSchemaParsed) Schema() (*ModelSchema, error) {
//	if m.DisplayName == "" {
//		return nil, fmt.Errorf("model has no display name")
//	}
//	if m.ModelType == "" {
//		return nil, fmt.Errorf("model has no model type: %s", m.DisplayName)
//	}
//	if m.CollectionName == "" {
//		return nil, fmt.Errorf("model has no collection name: %s", m.DisplayName)
//	}
//	if m.SingularName == "" {
//		return nil, fmt.Errorf("model has no singular name: %s", m.DisplayName)
//	}
//	if m.PluralName == "" {
//		return nil, fmt.Errorf("model has no plural name: %s", m.DisplayName)
//	}
//
//	var attrs []*AttributeSchema
//
//	for _, attr := range m.Attributes {
//		attrSchema, err := attr.Schema()
//		if err != nil {
//			return nil, err
//		}
//		attrs = append(attrs, attrSchema)
//	}
//
//	return &ModelSchema{
//		ModelType:      m.ModelType,
//		CollectionName: m.CollectionName,
//		DisplayName:    m.DisplayName,
//		SingularName:   m.SingularName,
//		PluralName:     m.PluralName,
//		Description:    m.Description,
//		Attributes:     attrs,
//	}, nil
//}
//
//type AttributeSchemaParsed struct {
//	Name               string `yaml:"name" json:"name"`
//	SimplifiedDataType string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
//	DetailedDataType   string `yaml:"detailedDataType" json:"detailedDataType"`
//
//	ShownInTable *bool  `yaml:"shownInTable" json:"shownInTable"`
//	Required     *bool  `yaml:"required" json:"required"`
//	Max          *int64 `yaml:"max" json:"max"`
//	Min          *int64 `yaml:"min" json:"min"`
//	MaxLength    *int   `yaml:"maxLength" json:"maxLength"`
//	MinLength    *int   `yaml:"minLength" json:"minLength"`
//	Private      *bool  `yaml:"private" json:"private"`
//	Editable     *bool  `yaml:"editable" json:"editable"`
//
//	Default  *string `yaml:"default" json:"default"`
//	Nullable *bool   `yaml:"nullable" json:"nullable"`
//	Unique   *bool   `yaml:"unique" json:"unique"`
//}
//
//func (a AttributeSchemaParsed) Schema() (*AttributeSchema, error) {
//	if a.Name == "" {
//		return nil, fmt.Errorf("attribute has no name")
//	}
//	if a.SimplifiedDataType == "" {
//		return nil, fmt.Errorf("attribute has no simple type: %s", a.Name)
//	}
//	if a.DetailedDataType == "" {
//		return nil, fmt.Errorf("attribute has no detailed type: %s", a.Name)
//	}
//
//	return &AttributeSchema{
//		Name:               a.Name,
//		SimplifiedDataType: a.SimplifiedDataType,
//		DetailedDataType:   a.DetailedDataType,
//
//		ShownInTable: checkDefault(true, a.ShownInTable),
//		Required:     checkDefault(false, a.Required),
//		Max:          checkDefault(int64(2147483647), a.Max),
//		Min:          checkDefault(int64(-2147483648), a.Min),
//		MaxLength:    checkDefault(-1, a.MaxLength),
//		MinLength:    checkDefault(-1, a.MinLength),
//		Private:      checkDefault(false, a.Private),
//		Editable:     checkDefault(true, a.Editable),
//
//		Default:  checkDefault("", a.Default),
//		Nullable: checkDefault(true, a.Nullable),
//		Unique:   checkDefault(false, a.Unique),
//	}, nil
//}
//
//func checkDefault[T any](defaultValue T, value *T) T {
//	if value == nil {
//		return defaultValue
//	}
//	return *value
//}
