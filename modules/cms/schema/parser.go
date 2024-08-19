package schema

import (
	"encoding/json"
	"fmt"
	"github.com/cosys-io/cosys/common"
	"io"
)

// ParseSchema parses json from a reader and scans it into a ModelSchema.
func ParseSchema(schema *ModelSchema, reader io.Reader) error {
	schemaParseable := modelParseable{}

	if err := json.NewDecoder(reader).Decode(&schemaParseable); err != nil {
		return err
	}

	parsedSchema, err := schemaParseable.Schema()
	if err != nil {
		return err
	}

	*schema = *parsedSchema
	return nil
}

// modelParseable is used to parse model schema json with default values.
// Fields that are missing in the json are parsed as nil pointers instead of zero values,
// which are then replaced with their respective default values.
type modelParseable struct {
	ModelType      string           `yaml:"modelType" json:"modelType"`
	CollectionName string           `yaml:"collectionName" json:"collectionName"`
	DisplayName    string           `yaml:"displayName" json:"displayName"`
	SingularName   string           `yaml:"singularName" json:"singularName"`
	PluralName     string           `yaml:"pluralName" json:"pluralName"`
	Description    string           `yaml:"description" json:"description"`
	Attributes     []*attrParseable `yaml:"attributes" json:"attributes"`
}

// Schema returns the corresponding ModelSchema with default values.
func (m modelParseable) Schema() (*ModelSchema, error) {
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

	var attrs []common.AttributeSchema
	var index int

	hasId := false
	for _, attr := range m.Attributes {
		if attr.Name == "id" {
			hasId = true
		}
	}
	if hasId {
		attrs = make([]common.AttributeSchema, len(m.Attributes))
		index = 0
	} else {
		attrs = make([]common.AttributeSchema, len(m.Attributes)+1)
		attrs[0] = &IdSchema
		index = 1
	}

	for _, attr := range m.Attributes {
		attrSchema, err := attr.Schema()
		if err != nil {
			return nil, err
		}
		attrs[index] = attrSchema
		index += 1
	}

	return &ModelSchema{
		modelType:      m.ModelType,
		collectionName: m.CollectionName,
		displayName:    m.DisplayName,
		singularName:   m.SingularName,
		pluralName:     m.PluralName,
		description:    m.Description,
		attributes:     attrs,
	}, nil
}

// attrParseable is used to parse attribute schema json with default values.
// Fields that are missing in the json are parsed as nil pointers instead of zero values,
// which are then replaced with their respective default values.
type attrParseable struct {
	Name               string `yaml:"name" json:"name"`
	SimplifiedDataType string `yaml:"simplifiedDataType" json:"simplifiedDataType"`
	DetailedDataType   string `yaml:"detailedDataType" json:"detailedDataType"`

	ShownInTable *bool     `yaml:"shownInTable" json:"shownInTable"`
	Required     *bool     `yaml:"required" json:"required"`
	Max          *int64    `yaml:"max" json:"max"`
	Min          *int64    `yaml:"min" json:"min"`
	MaxLength    *int      `yaml:"maxLength" json:"maxLength"`
	MinLength    *int      `yaml:"minLength" json:"minLength"`
	Private      *bool     `yaml:"private" json:"private"`
	Editable     *bool     `yaml:"editable" json:"editable"`
	Enum         *[]string `yaml:"enum" json:"enum"`

	Default  *string `yaml:"default" json:"default"`
	Nullable *bool   `yaml:"nullable" json:"nullable"`
	Unique   *bool   `yaml:"unique" json:"unique"`
}

// Schema returns the corresponding AttributeSchema with default values.
func (a attrParseable) Schema() (*AttributeSchema, error) {
	if a.Name == "" {
		return nil, fmt.Errorf("attribute has no name")
	}
	if a.SimplifiedDataType == "" {
		return nil, fmt.Errorf("attribute has no simple type: %s", a.Name)
	}
	if a.DetailedDataType == "" {
		return nil, fmt.Errorf("attribute has no detailed type: %s", a.Name)
	}

	return &AttributeSchema{
		name:               a.Name,
		simplifiedDataType: a.SimplifiedDataType,
		detailedDataType:   a.DetailedDataType,

		shownInTable: checkDefault(true, a.ShownInTable),
		required:     checkDefault(false, a.Required),
		max:          checkDefault(int64(2147483647), a.Max),
		min:          checkDefault(int64(-2147483648), a.Min),
		maxLength:    checkDefault(-1, a.MaxLength),
		minLength:    checkDefault(-1, a.MinLength),
		private:      checkDefault(false, a.Private),
		editable:     checkDefault(true, a.Editable),
		enum:         checkDefault([]string{}, a.Enum),

		defaultValue: checkDefault("", a.Default),
		nullable:     checkDefault(true, a.Nullable),
		unique:       checkDefault(false, a.Unique),
	}, nil
}

// checkDefault returns the value pointed to by a pointer, or a default value if the pointer is null.
func checkDefault[T any](defaultValue T, value *T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
