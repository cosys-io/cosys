package common

type BoolAttribute struct {
	*AttributeBase
}

func NewBoolAttribute(name, fieldName string) *BoolAttribute {
	return &BoolAttribute{
		&AttributeBase{
			name,
			fieldName,
			nil,
		},
	}
}

func (b BoolAttribute) Not() Condition {
	return &NestedCondition{
		Not,
		b,
		nil,
	}
}

func (b BoolAttribute) And(right Condition) Condition {
	return &NestedCondition{
		And,
		b,
		right,
	}
}

func (b BoolAttribute) Or(right Condition) Condition {
	return &NestedCondition{
		And,
		b,
		right,
	}
}
