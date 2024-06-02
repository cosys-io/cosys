package models

type AttributeSchema struct {
	Type string `yaml:"type"`

	Required        bool  `yaml:"required"`
	Max             int64 `yaml:"max"`
	Min             int64 `yaml:"min"`
	MaxLength       int   `yaml:"maxLength"`
	MinLength       int   `yaml:"minLength"`
	Private         bool  `yaml:"private"`
	NotConfigurable bool  `yaml:"notConfigurable"`

	Default     any  `yaml:"default"`
	NotNullable bool `yaml:"notNullable"`
	Unsigned    bool `yaml:"unsigned"`
	Unique      bool `yaml:"unique"`
}

func NewAttributeSchema(datatype string, options ...attrSchemaOption) *AttributeSchema {
	attributeSchema := &AttributeSchema{
		datatype,

		false,
		2147483647,
		-2147483648,
		-1,
		-1,
		false,
		false,

		nil,
		false,
		false,
		false,
	}

	for _, option := range options {
		option(attributeSchema)
	}

	return attributeSchema
}

type attrSchemaOption func(*AttributeSchema)

func Required(schema *AttributeSchema) {
	schema.Required = true
}

func Max(max int64) attrSchemaOption {
	return func(schema *AttributeSchema) {
		schema.Max = max
	}
}

func Min(min int64) attrSchemaOption {
	return func(schema *AttributeSchema) {
		schema.Min = min
	}
}

func MaxLength(maxLen int) attrSchemaOption {
	return func(schema *AttributeSchema) {
		schema.MaxLength = maxLen
	}
}

func MinLength(minLen int) attrSchemaOption {
	return func(schema *AttributeSchema) {
		schema.MinLength = minLen
	}
}

func Private(schema *AttributeSchema) {
	schema.Private = true
}

func NotConfigurable(schema *AttributeSchema) {
	schema.NotConfigurable = true
}

func Default(val any) attrSchemaOption {
	return func(schema *AttributeSchema) {
		schema.Default = val
	}
}

func NotNullable(schema *AttributeSchema) {
	schema.NotNullable = true
}

func Unsigned(schema *AttributeSchema) {
	schema.Unsigned = true
}

func Unique(schema *AttributeSchema) {
	schema.Unique = true
}
