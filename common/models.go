package common

// Model

type Model interface {
	New_() Entity
	All_() []Attribute
	Id_() *IntAttribute
	Name_() string
	Schema_() *ModelSchema
	Lifecycle_() Lifecycle
}

// Entity

type Entity interface {
}

// Attribute

type Attribute interface {
	Name() string
	FieldName() string

	Asc() *Order
	Desc() *Order

	Null() Condition
	NotNull() Condition
}

type AttributeBase struct {
	name      string
	fieldName string
	schema    *AttributeSchema
}

func (a AttributeBase) Name() string {
	return a.name
}

func (a AttributeBase) FieldName() string {
	return a.fieldName
}

func (a AttributeBase) Schema() *AttributeSchema {
	return a.schema
}

func (a AttributeBase) Asc() *Order {
	return &Order{
		a,
		Asc,
	}
}

func (a AttributeBase) Desc() *Order {
	return &Order{
		a,
		Desc,
	}
}

func (a AttributeBase) Null() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}

func (a AttributeBase) NotNull() Condition {
	return &ExpressionCondition{
		Null,
		a,
		nil,
	}
}

// Order

type Order struct {
	Attribute Attribute
	Order     OrderType
}

type OrderType string

const (
	Asc  OrderType = "Asc"
	Desc OrderType = "Desc"
)

// Condition

type Condition interface {
	Not() Condition
	And(Condition) Condition
	Or(Condition) Condition
}

type NestedCondition struct {
	Op    NestedOperation
	Left  Condition
	Right Condition
}

type NestedOperation string

const (
	Not NestedOperation = "Not"
	And NestedOperation = "And"
	Or  NestedOperation = "Or"
)

func (n NestedCondition) Not() Condition {
	return &NestedCondition{
		Not,
		n,
		nil,
	}
}

func (n NestedCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		n,
		right,
	}
}

func (n NestedCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
		n,
		right,
	}
}

type ExpressionCondition struct {
	Op    ExpressionOperation
	Left  Attribute
	Right any
}

type ExpressionOperation string

const (
	None ExpressionOperation = ""

	Eq    ExpressionOperation = "="
	Neq   ExpressionOperation = "<>"
	In    ExpressionOperation = "In"
	NotIn ExpressionOperation = "Not In"

	Lt  ExpressionOperation = "<"
	Gt  ExpressionOperation = ">"
	Lte ExpressionOperation = "<="
	Gte ExpressionOperation = ">="

	Contains    ExpressionOperation = "Contains"
	NotContains ExpressionOperation = "NotContains"
	StartsWith  ExpressionOperation = "StartsWith"
	EndsWith    ExpressionOperation = "EndsWith"

	Null    ExpressionOperation = "Is Null"
	NotNull ExpressionOperation = "Is Not Null"
)

func (e ExpressionCondition) Not() Condition {
	return &NestedCondition{
		Not,
		&e,
		nil,
	}
}

func (e ExpressionCondition) And(right Condition) Condition {
	return &NestedCondition{
		And,
		&e,
		right,
	}
}

func (e ExpressionCondition) Or(right Condition) Condition {
	return &NestedCondition{
		Or,
		&e,
		right,
	}
}

// IntAttribute

type IntAttribute struct {
	*AttributeBase
}

func NewIntAttribute(name, fieldName string) *IntAttribute {
	return &IntAttribute{
		&AttributeBase{
			name,
			fieldName,
			nil,
		},
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

// BoolAttribute

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

// StringAttribute

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
