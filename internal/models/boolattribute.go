package models

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

func (b *BoolAttribute) Not() Condition {
	return &NestedCondition{
		NOT,
		b,
		nil,
	}
}

func (b *BoolAttribute) And(right Condition) Condition {
	return &NestedCondition{
		AND,
		b,
		right,
	}
}

func (b *BoolAttribute) Or(right Condition) Condition {
	return &NestedCondition{
		AND,
		b,
		right,
	}
}
