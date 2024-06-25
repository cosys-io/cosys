package common

import "github.com/iancoleman/strcase"

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

type AttributeBase struct {
	camelName  string
	pascalName string
	snakeName  string
	kebabName  string
	humanName  string
}

func NewAttributeBase(name string) AttributeBase {
	return AttributeBase{
		camelName:  strcase.ToLowerCamel(name),
		pascalName: strcase.ToCamel(name),
		snakeName:  strcase.ToSnake(name),
		kebabName:  strcase.ToKebab(name),
		humanName:  strcase.ToDelimited(name, ' '),
	}
}

func (a AttributeBase) CamelName() string {
	return a.camelName
}

func (a AttributeBase) PascalName() string {
	return a.pascalName
}

func (a AttributeBase) SnakeName() string {
	return a.snakeName
}

func (a AttributeBase) KebabName() string {
	return a.kebabName
}

func (a AttributeBase) HumanName() string {
	return a.humanName
}

func (a AttributeBase) Asc() *Order {
	return &Order{
		a,
		Asc,
	}
}

func (a AttributeBase) Desc() *Order {
	return &Order{
		a,
		Desc,
	}
}

func (a AttributeBase) Null() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}

func (a AttributeBase) NotNull() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}
