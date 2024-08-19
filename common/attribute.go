package common

import "github.com/iancoleman/strcase"

// Attribute is a field of a model.
type Attribute interface {
	Asc() *Order
	Desc() *Order

	Null() Condition
	NotNull() Condition

	CamelName() string
	PascalName() string
	SnakeName() string
	KebabName() string
	HumanName() string
}

// attributeBase can be embedded into structs to provide them
// the methods to implement the Attribute interface.
type attributeBase struct {
	camelName  string
	pascalName string
	snakeName  string
	kebabName  string
	humanName  string
}

// newAttributeBase returns a new attributeBase
// based on the given attribute name.
func newAttributeBase(name string) attributeBase {
	return attributeBase{
		camelName:  strcase.ToLowerCamel(name),
		pascalName: strcase.ToCamel(name),
		snakeName:  strcase.ToSnake(name),
		kebabName:  strcase.ToKebab(name),
		humanName:  strcase.ToDelimited(name, ' '),
	}
}

func (a attributeBase) CamelName() string {
	return a.camelName
}

func (a attributeBase) PascalName() string {
	return a.pascalName
}

func (a attributeBase) SnakeName() string {
	return a.snakeName
}

func (a attributeBase) KebabName() string {
	return a.kebabName
}

func (a attributeBase) HumanName() string {
	return a.humanName
}

// Asc returns the ascending order-by condition for this attribute.
func (a attributeBase) Asc() *Order {
	return &Order{
		a,
		Asc,
	}
}

// Desc returns the descending order-by condition for this attribute.
func (a attributeBase) Desc() *Order {
	return &Order{
		a,
		Desc,
	}
}

// Null returns where condition, whether
// the value of the attribute is null.
func (a attributeBase) Null() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}

// NotNull returns where condition, whether
// the value of the attribute is not null.
func (a attributeBase) NotNull() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}
