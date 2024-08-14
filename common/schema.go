package common

// ModelSchema is a schema for a model.
type ModelSchema interface {
	CollectionName() string
	SingularName() string
	PluralName() string
	Attributes() []AttributeSchema
}

// AttributeSchema is a schema for a model attribute.
type AttributeSchema interface {
	Name() string
	SimplifiedDataType() string
	DetailedDataType() string

	Required() bool
	Max() int64
	Min() int64
	MaxLength() int
	MinLength() int
	Private() bool
	Editable() bool
	Enum() []string

	Default() string
	Nullable() bool
	Unique() bool
}

// modelSchema is an implementation of the ModelSchema interface.
type modelSchema struct {
	collectionName string
	singularName   string
	pluralName     string
	attributes     []AttributeSchema
}

func (s modelSchema) CollectionName() string {
	return s.collectionName
}

func (s modelSchema) SingularName() string {
	return s.singularName
}

func (s modelSchema) PluralName() string {
	return s.pluralName
}

func (s modelSchema) Attributes() []AttributeSchema {
	return s.attributes
}

// attributeSchema is an implementation of the AttributeSchema interface.
type attributeSchema struct {
	name               string
	simplifiedDataType string
	detailedDataType   string

	required  bool
	max       int64
	min       int64
	maxLength int
	minLength int
	private   bool
	editable  bool
	enum      []string

	defaultValue string
	nullable     bool
	unique       bool
}

func (s attributeSchema) Name() string {
	return s.name
}

func (s attributeSchema) SimplifiedDataType() string {
	return s.simplifiedDataType
}

func (s attributeSchema) DetailedDataType() string {
	return s.detailedDataType
}

func (s attributeSchema) Required() bool {
	return s.required
}

func (s attributeSchema) Max() int64 {
	return s.max
}

func (s attributeSchema) Min() int64 {
	return s.min
}

func (s attributeSchema) MaxLength() int {
	return s.maxLength
}

func (s attributeSchema) MinLength() int {
	return s.minLength
}

func (s attributeSchema) Private() bool {
	return s.private
}

func (s attributeSchema) Editable() bool {
	return s.editable
}

func (s attributeSchema) Enum() []string {
	return s.enum
}

func (s attributeSchema) Default() string {
	return s.defaultValue
}

func (s attributeSchema) Nullable() bool {
	return s.nullable
}

func (s attributeSchema) Unique() bool {
	return s.unique
}

// NewModelSchema returns a new model schema from the given names and attribute schemas.
func NewModelSchema(collection, singular, plural string, attrs ...AttributeSchema) ModelSchema {
	return &modelSchema{
		collectionName: collection,
		singularName:   singular,
		pluralName:     plural,
		attributes:     attrs,
	}
}

// NewAttrSchema returns a new attribute schema from the given name, types and configurations.
func NewAttrSchema(attrName, simpleType, detailedType string, opts ...AttrOption) AttributeSchema {
	schema := &attributeSchema{
		name:               attrName,
		simplifiedDataType: simpleType,
		detailedDataType:   detailedType,
		required:           false,
		max:                2147483647,
		min:                -2147483648,
		maxLength:          -1,
		minLength:          -1,
		private:            false,
		editable:           true,
		defaultValue:       "",
		nullable:           true,
		unique:             false,
	}

	for _, opt := range opts {
		opt(schema)
	}

	return schema
}

// AttrOption is a configuration for an attribute schema.
type AttrOption func(*attributeSchema)

// Required specifies that an attribute is required.
var Required AttrOption = func(schema *attributeSchema) {
	schema.required = true
}

// Max specifies that an attribute has the given maximum value.
func Max(max int64) AttrOption {
	return func(schema *attributeSchema) {
		schema.max = max
	}
}

// Min specifies that an attribute has the given minimum value.
func Min(min int64) AttrOption {
	return func(schema *attributeSchema) {
		schema.min = min
	}
}

// MaxLength specifies that an attribute has the given maximum length.
func MaxLength(max int) AttrOption {
	return func(schema *attributeSchema) {
		schema.maxLength = max
	}
}

// MinLength specifies that an attribute has the given minimum length.
func MinLength(min int) AttrOption {
	return func(schema *attributeSchema) {
		schema.minLength = min
	}
}

// Private specifies that an attribute is private.
var Private AttrOption = func(schema *attributeSchema) {
	schema.private = true
}

// NotEditable specifies that an attribute is not editable.
var NotEditable AttrOption = func(schema *attributeSchema) {
	schema.editable = false
}

// Enum specifies that an attribute is can only take the given values.
func Enum(enum []string) AttrOption {
	return func(schema *attributeSchema) {
		schema.enum = enum
	}
}

// Default specifies that an attribute has the given default value.
func Default(def string) AttrOption {
	return func(schema *attributeSchema) {
		schema.defaultValue = def
	}
}

// NotNullable specifies that an attribute is cannot be null.
var NotNullable AttrOption = func(schema *attributeSchema) {
	schema.nullable = false
}

// Unique specifies that an attribute is unique.
var Unique AttrOption = func(schema *attributeSchema) {
	schema.unique = true
}

// IdSchema is the schema for the id attribute.
var IdSchema = attributeSchema{
	name:               "id",
	simplifiedDataType: "Number",
	detailedDataType:   "Int",
	required:           false,
	max:                2147483647,
	min:                -2147483648,
	maxLength:          -1,
	minLength:          -1,
	private:            false,
	editable:           false,
	enum:               nil,
	nullable:           false,
	unique:             true,
}

// UuidSchema is the schema for the uuid attribute.
var UuidSchema = attributeSchema{
	name:               "uuid",
	simplifiedDataType: "String",
	detailedDataType:   "String",
	required:           false,
	max:                2147483647,
	min:                -2147483648,
	maxLength:          -1,
	minLength:          -1,
	private:            false,
	editable:           false,
	enum:               nil,
	nullable:           false,
	unique:             true,
}
