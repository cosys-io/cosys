package models

type Model interface {
	New_() Entity
	All_() []Attribute
	Id_() *IntAttribute
	Name_() string
}

type Entity interface {
}
