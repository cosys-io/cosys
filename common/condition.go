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

// NestedCondition is a where condition formed by performing the
// "not", "and", or "or" operations on other condition/s.
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

// Not returns the condition, formed by taking the
// logical negation of the nested condition.
func (n NestedCondition) Not() Condition {
	return &NestedCondition{
		Not,
		n,
		nil,
	}
}

// And returns the condition, formed by taking the
// logical conjunction of the nested condition and the given condition.
func (n NestedCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		n,
		right,
	}
}

// Or returns the condition, formed by taking the
// logical disjunction of the nested condition and the given condition.
func (n NestedCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
		n,
		right,
	}
}

// ExpressionCondition is a condition formed by
// performing operations a value.
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

// Not returns the condition, formed by taking the
// logical negation of the expression condition.
func (e ExpressionCondition) Not() Condition {
	return &NestedCondition{
		Not,
		&e,
		nil,
	}
}

// And returns the condition, formed by taking the
// logical conjunction of the expression condition and the given condition.
func (e ExpressionCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		&e,
		right,
	}
}

// Or returns the condition, formed by taking the
// logical disjunction of the expression condition and the given condition.
func (e ExpressionCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
		&e,
		right,
	}
}
