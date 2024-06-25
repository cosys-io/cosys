package common

type StringAttribute struct {
	*AttributeBase
}

func NewStringAttribute(name string) *StringAttribute {
	base := NewAttributeBase(name)

	return &StringAttribute{
		&base,
	}
}

func (s StringAttribute) Eq(right string) Condition {
	return &ExpressionCondition{
		Eq,
		s,
		right,
	}
}

func (s StringAttribute) NEq(right string) Condition {
	return &ExpressionCondition{
		Neq,
		s,
		right,
	}
}

func (s StringAttribute) In(right []string) Condition {
	return &ExpressionCondition{
		In,
		s,
		right,
	}
}

func (s StringAttribute) NotIn(right []string) Condition {
	return &ExpressionCondition{
		NotIn,
		s,
		right,
	}
}

func (s StringAttribute) Contains(right string) Condition {
	return &ExpressionCondition{
		Contains,
		s,
		right,
	}
}

func (s StringAttribute) NotContains(right string) Condition {
	return &ExpressionCondition{
		NotContains,
		s,
		right,
	}
}

func (s StringAttribute) StartsWith(right string) Condition {
	return &ExpressionCondition{
		StartsWith,
		s,
		right,
	}
}

func (s StringAttribute) EndsWith(right string) Condition {
	return &ExpressionCondition{
		EndsWith,
		s,
		right,
	}
}
