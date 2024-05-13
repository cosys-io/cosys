package models

type IAttribute interface {
	DBName() string
	StructName() string

	Asc() *Order
	Desc() *Order

	Null() Condition
	NotNull() Condition
}

type Attribute struct {
	dbname     string
	structname string
}

func (a *Attribute) DBName() string {
	return a.dbname
}

func (a *Attribute) StructName() string {
	return a.structname
}

func (a *Attribute) Asc() *Order {
	return &Order{
		a,
		ASC,
	}
}

func (a *Attribute) Desc() *Order {
	return &Order{
		a,
		DESC,
	}
}

func (a *Attribute) Null() Condition {
	return &ExpressionCondition{
		NULL,
		a,
		nil,
	}
}

func (a *Attribute) NotNull() Condition {
	return &ExpressionCondition{
		NULL,
		a,
		nil,
	}
}

// Order

type Order struct {
	Attribute IAttribute
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
	Left  IAttribute
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
