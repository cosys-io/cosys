package models

type Model interface {
	Model_Name() string
	Model_New() Entity
	Model_All() []IAttribute
	Model_Id() *IntAttribute
}

type Entity interface {
}
