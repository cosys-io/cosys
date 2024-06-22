package common

type IntAttribute struct {
	*AttributeBase
}

func NewIntAttribute(name string) *IntAttribute {
	base := NewAttributeBase(name)

	return &IntAttribute{
		&base,
	}
}

func (s IntAttribute) Eq(right int) Condition {
	return &ExpressionCondition{
		Eq,
		s,
		right,
	}
}

func (s IntAttribute) NEq(right int) Condition {
	return &ExpressionCondition{
		Neq,
		s,
		right,
	}
}

func (s IntAttribute) In(right []int) Condition {
	return &ExpressionCondition{
		In,
		s,
		right,
	}
}

func (s IntAttribute) NotIn(right []int) Condition {
	return &ExpressionCondition{
		NotIn,
		s,
		right,
	}
}

func (s IntAttribute) Lt(right int) Condition {
	return &ExpressionCondition{
		Lt,
		s,
		right,
	}
}

func (s IntAttribute) Gt(right int) Condition {
	return &ExpressionCondition{
		Gt,
		s,
		right,
	}
}

func (s IntAttribute) Lte(right int) Condition {
	return &ExpressionCondition{
		Lte,
		s,
		right,
	}
}

func (s IntAttribute) Gte(right int) Condition {
	return &ExpressionCondition{
		Gte,
		s,
		right,
	}
}
