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

type NestedCondition struct {
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

func (n NestedCondition) Not() Condition {
	return &NestedCondition{
		Not,
		n,
		nil,
	}
}

func (n NestedCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		n,
		right,
	}
}

func (n NestedCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
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

func (e ExpressionCondition) Not() Condition {
	return &NestedCondition{
		Not,
		&e,
		nil,
	}
}

func (e ExpressionCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		&e,
		right,
	}
}

func (e ExpressionCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
		&e,
		right,
	}
}
