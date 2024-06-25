package common

import "github.com/iancoleman/strcase"

type Model interface {
	New_() Entity
	All_() []Attribute
	Id_() *IntAttribute
	Schema_() *ModelSchema
	Lifecycle_() Lifecycle

	DBName_() string
	DisplayName_() string
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

type ModelBase struct {
	dbName             string
	displayName        string
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

func NewModelBase(dbName, displayName, singularName, pluralName string) ModelBase {
	return ModelBase{
		dbName:             dbName,
		displayName:        displayName,
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
}

func (m ModelBase) DBName_() string {
	return m.dbName
}

func (m ModelBase) DisplayName_() string {
	return m.displayName
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

type Entity interface {
}
