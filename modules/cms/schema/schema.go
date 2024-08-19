package schema

import "github.com/cosys-io/cosys/common"

// ModelSchema is an implementation of the ModelSchema common interface,
// with added getter methods for display name and description.
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

// AttributeSchema is an implementation of the AttributeSchema common interface,
// with the added getter method for whether the attribute is shown in table.
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

// NewModelSchema returns a new ModelSchema from the given names, descriptions and attribute schemas.
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

// NewAttrSchema returns a new AttributeSchema from the given names, types and configurations.
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

// AttrOption is a configuration for an attribute schema.
type AttrOption func(*AttributeSchema)

// Required specifies that the attribute is required.
var Required AttrOption = func(schema *AttributeSchema) {
	schema.required = true
}

// HideInTable specifies that the attribute is hidden in the table.
var HideInTable AttrOption = func(schema *AttributeSchema) {
	schema.shownInTable = false
}

// Max specifies that the attribute has the given maximum value.
func Max(max int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.max = max
	}
}

// Min specifies that the attribute has the given minimum value.
func Min(min int64) AttrOption {
	return func(schema *AttributeSchema) {
		schema.min = min
	}
}

// MaxLength specifies that the attribute has the given maximum length.
func MaxLength(max int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.maxLength = max
	}
}

// MinLength specifies that the attribute has the given minimum length.
func MinLength(min int) AttrOption {
	return func(schema *AttributeSchema) {
		schema.minLength = min
	}
}

// Private specifies that the attribute is private.
var Private AttrOption = func(schema *AttributeSchema) {
	schema.private = true
}

// NotEditable specifies that the attribute is not editable.
var NotEditable AttrOption = func(schema *AttributeSchema) {
	schema.editable = false
}

// Enum specifies that an attribute is can only take the given values.
func Enum(enum ...string) AttrOption {
	return func(schema *AttributeSchema) {
		schema.enum = enum
	}
}

// Default specifies that an attribute has the given default value.
func Default(def string) AttrOption {
	return func(schema *AttributeSchema) {
		schema.defaultValue = def
	}
}

// NotNullable specifies that an attribute is cannot be null.
var NotNullable AttrOption = func(schema *AttributeSchema) {
	schema.nullable = false
}

// Unique specifies that an attribute is unique.
var Unique AttrOption = func(schema *AttributeSchema) {
	schema.unique = true
}

// IdSchema is the schema for the id attribute.
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

// UuidSchema is the schema for the uuid attribute.
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
