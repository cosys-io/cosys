package models

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

func (a *AttributeBase) Name() string {
	return a.name
}

func (a *AttributeBase) FieldName() string {
	return a.fieldName
}

func (a *AttributeBase) Schema() *AttributeSchema {
	return a.schema
}

func (a *AttributeBase) Asc() *Order {
	return &Order{
		a,
		ASC,
	}
}

func (a *AttributeBase) Desc() *Order {
	return &Order{
		a,
		DESC,
	}
}

func (a *AttributeBase) Null() Condition {
	return &ExpressionCondition{
		NULL,
		a,
		nil,
	}
}

func (a *AttributeBase) NotNull() Condition {
	return &ExpressionCondition{
		NULL,
		a,
		nil,
	}
}

// Order

type Order struct {
	Attribute Attribute
	Order     OrderType
}

type OrderType string

const (
	ASC  OrderType = "ASC"
	DESC OrderType = "DESC"
)

// Condition

type Condition interface {
	Not() Condition
	And(Condition) Condition
	Or(Condition) Condition
}

type NestedCondition struct {
	Op    NestedOperation
	Left  Condition
	Right Condition
}

type NestedOperation string

const (
	NOT NestedOperation = "NOT"
	AND NestedOperation = "AND"
	OR  NestedOperation = "OR"
)

func (n *NestedCondition) Not() Condition {
	return &NestedCondition{
		NOT,
		n,
		nil,
	}
}

func (n *NestedCondition) And(right Condition) Condition {
	return &NestedCondition{
		AND,
		n,
		right,
	}
}

func (n *NestedCondition) Or(right Condition) Condition {
	return &NestedCondition{
		OR,
		n,
		right,
	}
}

type ExpressionCondition struct {
	Op    ExpressionOperation
	Left  Attribute
	Right any
}

type ExpressionOperation string

const (
	NONE ExpressionOperation = ""

	EQ    ExpressionOperation = "="
	NEQ   ExpressionOperation = "<>"
	IN    ExpressionOperation = "IN"
	NOTIN ExpressionOperation = "NOT IN"

	LT  ExpressionOperation = "<"
	GT  ExpressionOperation = ">"
	LTE ExpressionOperation = "<="
	GTE ExpressionOperation = ">="

	CONTAINS    ExpressionOperation = "CONTAINS"
	NOTCONTAINS ExpressionOperation = "NOTCONTAINS"
	STARTSWITH  ExpressionOperation = "STARTSWITH"
	ENDSWITH    ExpressionOperation = "ENDSWITH"

	NULL    ExpressionOperation = "IS NULL"
	NOTNULL ExpressionOperation = "IS NOT NULL"
)

func (e ExpressionCondition) Not() Condition {
	return &NestedCondition{
		NOT,
		&e,
		nil,
	}
}

func (e ExpressionCondition) And(right Condition) Condition {
	return &NestedCondition{
		AND,
		&e,
		right,
	}
}

func (e ExpressionCondition) Or(right Condition) Condition {
	return &NestedCondition{
		OR,
		&e,
		right,
	}
}
