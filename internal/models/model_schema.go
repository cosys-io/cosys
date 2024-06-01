package models

type ModelSchema struct {
	ModelType__      string                      `yaml:"modelType"`
	CollectionName__ string                      `yaml:"collectionName"`
	Info__           *Info                       `yaml:"info"`
	Attributes__     map[string]*AttributeSchema `yaml:"attributes"`
}

type Info struct {
	DisplayName  string `yaml:"displayName"`
	SingularName string `yaml:"singularName"`
	PluralName   string `yaml:"pluralName"`
	Description  string `yaml:"description"`
}

func NewModelSchema(modelType, collection, display, singular, plural, description string) *ModelSchema {
	return &ModelSchema{
		modelType,
		collection,
		&Info{
			display,
			singular,
			plural,
			description,
		},
		map[string]*AttributeSchema{},
	}
}

func EmptyModelSchema() *ModelSchema {
	return &ModelSchema{
		"",
		"",
		&Info{
			"",
			"",
			"",
			"",
		},
		map[string]*AttributeSchema{},
	}
}

func (m *ModelSchema) CollectionName_() string {
	return m.CollectionName__
}

func (m *ModelSchema) DisplayName() string {
	return m.Info__.DisplayName
}

func (m *ModelSchema) SingularName() string {
	return m.Info__.SingularName
}

func (m *ModelSchema) PluralName() string {
	return m.Info__.PluralName
}

func (m *ModelSchema) Description() string {
	return m.Info__.Description
}
