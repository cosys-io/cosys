package common

// Order

type Order struct {
	Attribute Attribute
	Order     OrderType
}

type OrderType string

const (
	Asc  OrderType = "Asc"
	Desc OrderType = "Desc"
)

// Condition

type Condition interface {
	Not() Condition
	And(Condition) Condition
	Or(Condition) Condition
}

type nestedCondition struct {
	Op    NestedOperation
	Left  Condition
	Right Condition
}

type NestedOperation string

const (
	Not NestedOperation = "Not"
	And NestedOperation = "And"
	Or  NestedOperation = "Or"
)

func (n nestedCondition) Not() Condition {
	return &nestedCondition{
		Not,
		n,
		nil,
	}
}

func (n nestedCondition) And(right Condition) Condition {
	return &nestedCondition{
		And,
		n,
		right,
	}
}

func (n nestedCondition) Or(right Condition) Condition {
	return &nestedCondition{
		Or,
		n,
		right,
	}
}

type expressionCondition struct {
	Op    ExpressionOperation
	Left  Attribute
	Right any
}

type ExpressionOperation string

const (
	None ExpressionOperation = ""

	Eq    ExpressionOperation = "="
	Neq   ExpressionOperation = "<>"
	In    ExpressionOperation = "In"
	NotIn ExpressionOperation = "Not In"

	Lt  ExpressionOperation = "<"
	Gt  ExpressionOperation = ">"
	Lte ExpressionOperation = "<="
	Gte ExpressionOperation = ">="

	Contains    ExpressionOperation = "Contains"
	NotContains ExpressionOperation = "NotContains"
	StartsWith  ExpressionOperation = "StartsWith"
	EndsWith    ExpressionOperation = "EndsWith"

	Null    ExpressionOperation = "Is Null"
	NotNull ExpressionOperation = "Is Not Null"
)

func (e expressionCondition) Not() Condition {
	return &nestedCondition{
		Not,
		&e,
		nil,
	}
}

func (e expressionCondition) And(right Condition) Condition {
	return &nestedCondition{
		And,
		&e,
		right,
	}
}

func (e expressionCondition) Or(right Condition) Condition {
	return &nestedCondition{
		Or,
		&e,
		right,
	}
}
