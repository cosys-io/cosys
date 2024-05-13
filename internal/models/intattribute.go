package models

type IntAttribute struct {
	*Attribute
}

func NewIntAttribute(dbname, structname string) *IntAttribute {
	return &IntAttribute{
		&Attribute{
			dbname,
			structname,
		},
	}
}

func (s *IntAttribute) Eq(right int) Condition {
	return &ExpressionCondition{
		EQ,
		s,
		right,
	}
}

func (s *IntAttribute) NEq(right int) Condition {
	return &ExpressionCondition{
		NEQ,
		s,
		right,
	}
}

func (s *IntAttribute) In(right []int) Condition {
	return &ExpressionCondition{
		IN,
		s,
		right,
	}
}

func (s *IntAttribute) NotIn(right []int) Condition {
	return &ExpressionCondition{
		NOTIN,
		s,
		right,
	}
}

func (s *IntAttribute) Lt(right int) Condition {
	return &ExpressionCondition{
		LT,
		s,
		right,
	}
}

func (s *IntAttribute) Gt(right int) Condition {
	return &ExpressionCondition{
		GT,
		s,
		right,
	}
}

func (s *IntAttribute) Lte(right int) Condition {
	return &ExpressionCondition{
		LTE,
		s,
		right,
	}
}
func (s *IntAttribute) Gte(right int) Condition {
	return &ExpressionCondition{
		GTE,
		s,
		right,
	}
}
