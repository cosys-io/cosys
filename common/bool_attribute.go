package common

// BoolAttribute is an attribute of boolean datatype.
type BoolAttribute struct {
	*attributeBase
}

// NewBoolAttribute returns a new boolean attribute with the given name.
func NewBoolAttribute(name string) BoolAttribute {
	base := newAttributeBase(name)

	return BoolAttribute{
		&base,
	}
}

// Not returns where condition, whether the value of
// the boolean attribute is false.
func (b BoolAttribute) Not() Condition {
	return &nestedCondition{
		Not,
		b,
		nil,
	}
}

// And returns where condition, whether the value of
// the boolean attribute and the given condition are true.
func (b BoolAttribute) And(right Condition) Condition {
	return &nestedCondition{
		And,
		b,
		right,
	}
}

// Or returns where condition, whether the value of
// the boolean attribute or the given condition is true.
func (b BoolAttribute) Or(right Condition) Condition {
	return &nestedCondition{
		And,
		b,
		right,
	}
}
