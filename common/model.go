package common

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"reflect"
)

// Model is a content model.
type Model interface {
	New_() Entity
	Attributes_() []Attribute
	IdAttribute_() Attribute
	Schema_() *ModelSchema

	GetLifecycleHook_(event string, uid string) (LifecycleHook, error)
	CallLifecycle_(event string, query EventQuery) error
	AddLifecycleHook_(event string, hook LifecycleHook) (string, error)
	UpdateLifecycleHook_(event string, uid string, hook LifecycleHook) error
	RemoveLifecycleHook_(event string, uid string) error

	DBName_() string
	SingularCamelName_() string
	PluralCamelName_() string
	SingularPascalName_() string
	PluralPascalName_() string
	SingularSnakeName_() string
	PluralSnakeName_() string
	SingularKebabName_() string
	PluralKebabName_() string
	SingularHumanName_() string
	PluralHumanName_() string
}

// NewModel returns a new model that is built upon ModelBase.
func NewModel[E Entity, M Model](dbname, singularName, pluralName string) (M, error) {
	base := &ModelBase{
		entity: *new(E),

		lifecycle: NewLifecycle(),

		dbname:             dbname,
		singularCamelName:  strcase.ToLowerCamel(singularName),
		pluralCamelName:    strcase.ToLowerCamel(pluralName),
		singularPascalName: strcase.ToCamel(singularName),
		pluralPascalName:   strcase.ToCamel(pluralName),
		singularSnakeName:  strcase.ToSnake(singularName),
		pluralSnakeName:    strcase.ToSnake(pluralName),
		singularKebabName:  strcase.ToKebab(singularName),
		pluralKebabName:    strcase.ToKebab(pluralName),
		singularHumanName:  strcase.ToDelimited(singularName, ' '),
		pluralHumanName:    strcase.ToDelimited(pluralName, ' '),
	}

	model := new(M)
	if reflect.ValueOf(model).Elem().Kind() != reflect.Struct {
		return *new(M), fmt.Errorf("model is not a struct")
	}

	modelValue := reflect.ValueOf(model).Elem()
	foundId := false
	for i := range modelValue.NumField() {
		name := modelValue.Type().Field(i).Name
		fieldValue := modelValue.Field(i)
		field := fieldValue.Interface()

		var attr Attribute
		switch field.(type) {
		case BoolAttribute:
			attr = NewBoolAttribute(name)
		case IntAttribute:
			attr = NewIntAttribute(name)
			if name == "Id" {
				foundId = true
				base.idAttribute = attr
			}
		case StringAttribute:
			attr = NewStringAttribute(name)
		case *ModelBase:
			fieldValue.Set(reflect.ValueOf(base))
			continue
		default:
			continue
		}
		fieldValue.Set(reflect.ValueOf(attr))
		base.attributes = append(base.attributes, attr)
	}
	if !foundId {
		return *new(M), fmt.Errorf("id attribute not found in model %s", pluralName)
	}

	return *model, nil
}

// ModelBase pointers can be embedded into structs to provide them with
// the methods to implement the Model interface.
type ModelBase struct {
	entity      Entity
	idAttribute Attribute
	attributes  []Attribute

	lifecycle Lifecycle

	dbname             string
	singularCamelName  string
	pluralCamelName    string
	singularPascalName string
	pluralPascalName   string
	singularSnakeName  string
	pluralSnakeName    string
	singularKebabName  string
	pluralKebabName    string
	singularHumanName  string
	pluralHumanName    string
}

// New_ returns a new entity of the model.
func (m ModelBase) New_() Entity {
	return reflect.Zero(reflect.TypeOf(m.entity)).Interface().(Entity)
}

// IdAttribute_ returns the id attribute of the model.
func (m ModelBase) IdAttribute_() Attribute {
	return m.idAttribute
}

// Attributes_ return a slice of all attributes of the model.
func (m ModelBase) Attributes_() []Attribute {
	return m.attributes
}

// GetLifecycleHook_ returns a hook specified by the given uid
// for the given event in the lifecycle of the model.
func (m ModelBase) GetLifecycleHook_(event string, uid string) (LifecycleHook, error) {
	return m.lifecycle.Get(event, uid)
}

// CallLifecycle_ calls all hooks for the given event
// in the lifecycle of the model.
func (m ModelBase) CallLifecycle_(event string, query EventQuery) error {
	return m.lifecycle.Call(event, query)
}

// AddLifecycleHook_ adds a hook for the given event to the lifecycle of the model,
// and returns a generated uid for updating or removing.
func (m ModelBase) AddLifecycleHook_(event string, hook LifecycleHook) (string, error) {
	return m.lifecycle.Add(event, hook)
}

// UpdateLifecycleHook_ updates a hook specified by the given uid
// for the given event in the lifecycle of the model.
func (m ModelBase) UpdateLifecycleHook_(event string, uid string, hook LifecycleHook) error {
	return m.lifecycle.Update(event, uid, hook)
}

// RemoveLifecycleHook_ removes a hook specified by the given uid
// for the given event from the lifecycle of the model.
func (m ModelBase) RemoveLifecycleHook_(event string, uid string) error {
	return m.lifecycle.Remove(event, uid)
}

func (m ModelBase) DBName_() string {
	return m.dbname
}

func (m ModelBase) SingularCamelName_() string {
	return m.singularCamelName
}

func (m ModelBase) PluralCamelName_() string {
	return m.pluralCamelName
}

func (m ModelBase) SingularPascalName_() string {
	return m.singularPascalName
}

func (m ModelBase) PluralPascalName_() string {
	return m.pluralPascalName
}

func (m ModelBase) SingularSnakeName_() string {
	return m.singularSnakeName
}

func (m ModelBase) PluralSnakeName_() string {
	return m.pluralSnakeName
}

func (m ModelBase) SingularKebabName_() string {
	return m.singularKebabName
}

func (m ModelBase) PluralKebabName_() string {
	return m.pluralKebabName
}

func (m ModelBase) SingularHumanName_() string {
	return m.singularHumanName
}

func (m ModelBase) PluralHumanName_() string {
	return m.pluralHumanName
}

// Entity is a record/document.
type Entity interface {
}
