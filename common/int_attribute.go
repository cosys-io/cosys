package common

// IntAttribute is an attribute of integer datatype.
type IntAttribute struct {
	*attributeBase
}

// NewIntAttribute returns a new integer attribute with the given name.
func NewIntAttribute(name string) IntAttribute {
	base := newAttributeBase(name)

	return IntAttribute{
		&base,
	}
}

// Eq returns the where condition, whether the value of
// the integer attribute is equals to the given integer.
func (s IntAttribute) Eq(right int) Condition {
	return &expressionCondition{
		Eq,
		s,
		right,
	}
}

// NEq returns the where condition, whether the value of
// the integer attribute is not equals to the given integer.
func (s IntAttribute) NEq(right int) Condition {
	return &expressionCondition{
		Neq,
		s,
		right,
	}
}

// In returns the where condition, whether the value of
// the integer attribute is in the given slice of integers.
func (s IntAttribute) In(right []int) Condition {
	return &expressionCondition{
		In,
		s,
		right,
	}
}

// NotIn returns the where condition, whether the value of
// the integer attribute is not in the given slice of integers.
func (s IntAttribute) NotIn(right []int) Condition {
	return &expressionCondition{
		NotIn,
		s,
		right,
	}
}

// Lt returns the where condition, whether the value of
// the integer attribute is less than the given integer.
func (s IntAttribute) Lt(right int) Condition {
	return &expressionCondition{
		Lt,
		s,
		right,
	}
}

// Gt returns the where condition, whether the value of
// the integer attribute is greater than the given integer.
func (s IntAttribute) Gt(right int) Condition {
	return &expressionCondition{
		Gt,
		s,
		right,
	}
}

// Lte returns the where condition, whether the value of
// the integer attribute is less than or equals to the given integer.
func (s IntAttribute) Lte(right int) Condition {
	return &expressionCondition{
		Lte,
		s,
		right,
	}
}

// Gte returns the where condition, whether the value of
// the integer attribute is greater than or equals to the given integer.
func (s IntAttribute) Gte(right int) Condition {
	return &expressionCondition{
		Gte,
		s,
		right,
	}
}
