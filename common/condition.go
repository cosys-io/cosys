package common

// Order

// Order is an order-by condition.
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

// Condition is a where condition.
type Condition interface {
	Not() Condition
	And(Condition) Condition
	Or(Condition) Condition
}

// nestedCondition is a where condition formed by performing the
// "not", "and", or "or" operations on other condition/s.
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

// Not returns the condition, formed by taking the
// logical negation of the nested condition.
func (n nestedCondition) Not() Condition {
	return &nestedCondition{
		Not,
		n,
		nil,
	}
}

// And returns the condition, formed by taking the
// logical conjunction of the nested condition and the given condition.
func (n nestedCondition) And(right Condition) Condition {
	return &nestedCondition{
		And,
		n,
		right,
	}
}

// Or returns the condition, formed by taking the
// logical disjunction of the nested condition and the given condition.
func (n nestedCondition) Or(right Condition) Condition {
	return &nestedCondition{
		Or,
		n,
		right,
	}
}

// expressionCondition is a condition formed by
// performing operations a value.
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

// Not returns the condition, formed by taking the
// logical negation of the expression condition.
func (e expressionCondition) Not() Condition {
	return &nestedCondition{
		Not,
		&e,
		nil,
	}
}

// And returns the condition, formed by taking the
// logical conjunction of the expression condition and the given condition.
func (e expressionCondition) And(right Condition) Condition {
	return &nestedCondition{
		And,
		&e,
		right,
	}
}

// Or returns the condition, formed by taking the
// logical disjunction of the expression condition and the given condition.
func (e expressionCondition) Or(right Condition) Condition {
	return &nestedCondition{
		Or,
		&e,
		right,
	}
}
