package models

type BoolAttribute struct {
	*Attribute
}

func NewBoolAttribute(dbname, structname string) *BoolAttribute {
	return &BoolAttribute{
		&Attribute{
			dbname,
			structname,
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
