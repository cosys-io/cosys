package schema

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
)

// ToModelSerializable returns a serializable model schema from the given schema.
func ToModelSerializable(schema common.ModelSchema) (ModelSerializable, error) {
	schemaCMS, ok := schema.(*ModelSchema)
	if !ok {
		return ModelSerializable{}, fmt.Errorf("model schema is not a CMS model schema: %s", schema.PluralName())
	}

	attrsSerializable := make([]*AttrSerializable, len(schemaCMS.attributes))
	for i, attr := range schemaCMS.attributes {
		attrCMS, ok := attr.(*AttributeSchema)
		if !ok {
			return ModelSerializable{}, fmt.Errorf("attribute schema is not a CMS attribute schema: %s", attr.Name())
		}
		attrsSerializable[i] = ToAttrSerializable(attrCMS)
	}

	return ModelSerializable{
		ModelType:      schemaCMS.modelType,
		CollectionName: schemaCMS.collectionName,
		DisplayName:    schemaCMS.displayName,
		SingularName:   schemaCMS.singularName,
		PluralName:     schemaCMS.pluralName,
		Description:    schemaCMS.description,
		Attributes:     attrsSerializable,
	}, nil
}

// ToAttrSerializable returns a serializable attribute schema from the given schema.
func ToAttrSerializable(schema *AttributeSchema) *AttrSerializable {
	return &AttrSerializable{
		Name:               schema.name,
		SimplifiedDataType: schema.simplifiedDataType,
		DetailedDataType:   schema.detailedDataType,

		ShownInTable: schema.shownInTable,
		Required:     schema.required,
		Max:          schema.max,
		Min:          schema.min,
		MaxLength:    schema.maxLength,
		MinLength:    schema.minLength,
		Private:      schema.private,
		Editable:     schema.editable,
		Enum:         schema.enum,

		DefaultValue: schema.defaultValue,
		Nullable:     schema.nullable,
		Unique:       schema.unique,
	}
}

// ModelSerializable is used to serialize model schema into json.
type ModelSerializable struct {
	ModelType      string              `json:"modelType"`
	CollectionName string              `json:"collectionName"`
	DisplayName    string              `json:"displayName"`
	SingularName   string              `json:"singularName"`
	PluralName     string              `json:"pluralName"`
	Description    string              `json:"description"`
	Attributes     []*AttrSerializable `json:"attributes"`
}

// AttrSerializable is used to serialize attribute schema into json.
type AttrSerializable struct {
	Name               string `json:"name"`
	SimplifiedDataType string `json:"simplifiedDataType"`
	DetailedDataType   string `json:"detailedDataType"`

	ShownInTable bool     `json:"shownInTable"`
	Required     bool     `json:"required"`
	Max          int64    `json:"max"`
	Min          int64    `json:"min"`
	MaxLength    int      `json:"maxLength"`
	MinLength    int      `json:"minLength"`
	Private      bool     `json:"private"`
	Editable     bool     `json:"editable"`
	Enum         []string `json:"enum"`

	DefaultValue string `json:"default"`
	Nullable     bool   `json:"nullable"`
	Unique       bool   `json:"unique"`
}
