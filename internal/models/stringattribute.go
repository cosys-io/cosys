package models

type StringAttribute struct {
	*AttributeBase
}

func NewStringAttribute(name, fieldName string) *StringAttribute {
	return &StringAttribute{
		&AttributeBase{
			name,
			fieldName,
			nil,
		},
	}
}

func (s *StringAttribute) Eq(right string) Condition {
	return &ExpressionCondition{
		EQ,
		s,
		right,
	}
}

func (s *StringAttribute) NEq(right string) Condition {
	return &ExpressionCondition{
		NEQ,
		s,
		right,
	}
}

func (s *StringAttribute) In(right []string) Condition {
	return &ExpressionCondition{
		IN,
		s,
		right,
	}
}

func (s *StringAttribute) NotIn(right []string) Condition {
	return &ExpressionCondition{
		NOTIN,
		s,
		right,
	}
}

func (s *StringAttribute) Contains(right string) Condition {
	return &ExpressionCondition{
		CONTAINS,
		s,
		right,
	}
}

func (s *StringAttribute) NotContains(right string) Condition {
	return &ExpressionCondition{
		NOTCONTAINS,
		s,
		right,
	}
}

func (s *StringAttribute) StartsWith(right string) Condition {
	return &ExpressionCondition{
		STARTSWITH,
		s,
		right,
	}
}

func (s *StringAttribute) EndsWith(right string) Condition {
	return &ExpressionCondition{
		ENDSWITH,
		s,
		right,
	}
}
