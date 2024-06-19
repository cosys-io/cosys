package common

// Model

type Model interface {
	New_() Entity
	All_() []Attribute
	Id_() *IntAttribute
	Name_() string
	Schema_() *ModelSchema
	Lifecycle_() Lifecycle
}

// Entity

type Entity interface {
}

// Attribute

type Attribute interface {
	Name() string
	FieldName() string

	Asc() *Order
	Desc() *Order

	Null() Condition
	NotNull() Condition
}

type AttributeBase struct {
	name      string
	fieldName string
	schema    *AttributeSchema
}

func (a AttributeBase) Name() string {
	return a.name
}

func (a AttributeBase) FieldName() string {
	return a.fieldName
}

func (a AttributeBase) Schema() *AttributeSchema {
	return a.schema
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
