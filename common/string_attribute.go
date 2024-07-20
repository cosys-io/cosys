package common

// StringAttribute is an attribute of string datatype.
type StringAttribute struct {
	*attributeBase
}

// NewStringAttribute returns a new string attribute with the given name.
func NewStringAttribute(name string) StringAttribute {
	base := newAttributeBase(name)

	return StringAttribute{
		&base,
	}
}

// Eq returns the where condition, whether the value of
// the string attribute is equals to the given string.
func (s StringAttribute) Eq(right string) Condition {
	return &expressionCondition{
		Eq,
		s,
		right,
	}
}

// NEq returns the where condition, whether the value of
// the string attribute is not equals to the given string.
func (s StringAttribute) NEq(right string) Condition {
	return &expressionCondition{
		Neq,
		s,
		right,
	}
}

// In returns the where condition, whether the value of
// the string attribute is in the given slice of strings.
func (s StringAttribute) In(right []string) Condition {
	return &expressionCondition{
		In,
		s,
		right,
	}
}

// NotIn returns the where condition, whether the value of
// the string attribute is in the given slice of strings.
func (s StringAttribute) NotIn(right []string) Condition {
	return &expressionCondition{
		NotIn,
		s,
		right,
	}
}

// Contains returns the where condition, whether the value of
// the string attribute contains the given string.
func (s StringAttribute) Contains(right string) Condition {
	return &expressionCondition{
		Contains,
		s,
		right,
	}
}

// NotContains returns the where condition, whether the value of
// the string attribute does not contain the given string.
func (s StringAttribute) NotContains(right string) Condition {
	return &expressionCondition{
		NotContains,
		s,
		right,
	}
}

// StartsWith returns the where condition, whether the value of
// the string attribute starts with the given string.
func (s StringAttribute) StartsWith(right string) Condition {
	return &expressionCondition{
		StartsWith,
		s,
		right,
	}
}

// EndsWith returns the where condition, whether the value of
// the string attribute ends with the given string.
func (s StringAttribute) EndsWith(right string) Condition {
	return &expressionCondition{
		EndsWith,
		s,
		right,
	}
}
