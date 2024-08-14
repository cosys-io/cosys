package schema

import "github.com/cosys-io/cosys/common"

type ModelSchema struct {
	modelType      string
	collectionName string
	displayName    string
	singularName   string
	pluralName     string
	description    string
	attributes     []common.AttributeSchema
}

func (m ModelSchema) ModelType() string {
	return m.modelType
}

func (m ModelSchema) CollectionName() string {
	return m.collectionName
}

func (m ModelSchema) DisplayName() string {
	return m.displayName
}

func (m ModelSchema) SingularName() string {
	return m.singularName
}

func (m ModelSchema) PluralName() string {
	return m.pluralName
}

func (m ModelSchema) Description() string {
	return m.description
}

func (m ModelSchema) Attributes() []common.AttributeSchema {
	return m.attributes
}

type AttributeSchema struct {
	name               string
	simplifiedDataType string
	detailedDataType   string

	shownInTable bool
	required     bool
	max          int64
	min          int64
	maxLength    int
	minLength    int
	private      bool
	editable     bool
	enum         []string

	defaultValue string
	nullable     bool
	unique       bool
}

func (a AttributeSchema) Name() string {
	return a.name
}

func (a AttributeSchema) SimplifiedDataType() string {
	return a.simplifiedDataType
}

func (a AttributeSchema) DetailedDataType() string {
	return a.detailedDataType
}

func (a AttributeSchema) ShownInTable() bool {
	return a.shownInTable
}

func (a AttributeSchema) Required() bool {
	return a.required
}

func (a AttributeSchema) Max() int64 {
	return a.max
}

func (a AttributeSchema) Min() int64 {
	return a.min
}

func (a AttributeSchema) MaxLength() int {
	return a.maxLength
}

func (a AttributeSchema) MinLength() int {
	return a.minLength
}

func (a AttributeSchema) Private() bool {
	return a.private
}

func (a AttributeSchema) Editable() bool {
	return a.editable
}

func (a AttributeSchema) Enum() []string {
	return a.enum
}

func (a AttributeSchema) Default() string {
	return a.defaultValue
}

func (a AttributeSchema) Nullable() bool {
	return a.nullable
}

func (a AttributeSchema) Unique() bool {
	return a.unique
}

func NewModelSchema(collection, display, singular, plural, description string, attrs ...*AttributeSchema) *ModelSchema {
	commonAttrs := make([]common.AttributeSchema, len(attrs))
	for i, attr := range attrs {
		commonAttrs[i] = attr
	}

	return &ModelSchema{
		modelType:      "collection",
		collectionName: collection,
		displayName:    display,
		singularName:   singular,
		pluralName:     plural,
		description:    description,
		attributes:     commonAttrs,
	}
}

func NewAttrSchema(attrName, simpleType, detailedType string, opts ...AttrOption) *AttributeSchema {
	schema := &AttributeSchema{
		name:               attrName,
		simplifiedDataType: simpleType,
		detailedDataType:   detailedType,
		required:           false,
		shownInTable:       true,
		max:                2147483647,
		min:                -2147483648,
		maxLength:          -1,
		minLength:          -1,
		private:            false,
		editable:           true,
		enum:               nil,
		defaultValue:       "",
		nullable:           true,
		unique:             false,
	}

	for _, opt := range opts {
		opt(schema)
	}

	return schema
}

type AttrOption func(*AttributeSchema)

var Required AttrOption = func(schema *AttributeSchema) {
	schema.required = true
}

var HideInTable AttrOption = func(schema *AttributeSchema) {
	schema.shownInTable = false
}

func Max(max int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.max = max
	}
}

func Min(min int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.min = min
	}
}

func MaxLength(max int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.maxLength = max
	}
}

func MinLength(min int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.minLength = min
	}
}

var Private AttrOption = func(schema *AttributeSchema) {
	schema.private = true
}

var NotEditable AttrOption = func(schema *AttributeSchema) {
	schema.editable = false
}

func Enum(enum ...string) AttrOption {
	return func(schema *AttributeSchema) {
		schema.enum = enum
	}
}

func Default(def string) AttrOption {
	return func(schema *AttributeSchema) {
		schema.defaultValue = def
	}
}

var NotNullable AttrOption = func(schema *AttributeSchema) {
	schema.nullable = false
}

var Unique AttrOption = func(schema *AttributeSchema) {
	schema.unique = true
}

var IdSchema = AttributeSchema{
	name:               "id",
	simplifiedDataType: "Number",
	detailedDataType:   "Int",

	shownInTable: true,
	required:     false,
	max:          2147483647,
	min:          -2147483648,
	maxLength:    -1,
	minLength:    -1,
	private:      false,
	editable:     false,
	enum:         nil,

	defaultValue: "",
	nullable:     false,
	unique:       true,
}

var UuidSchema = AttributeSchema{
	name:               "uuid",
	simplifiedDataType: "String",
	detailedDataType:   "String",

	shownInTable: true,
	required:     false,
	max:          2147483647,
	min:          -2147483648,
	maxLength:    -1,
	minLength:    -1,
	private:      false,
	editable:     false,
	enum:         nil,

	defaultValue: "",
	nullable:     false,
	unique:       true,
}
