package common

type ModelSchema struct {
	ModelType      string                      `yaml:"modelType"`
	CollectionName string                      `yaml:"collectionName"`
	DisplayName    string                      `yaml:"displayName"`
	SingularName   string                      `yaml:"singularName"`
	PluralName     string                      `yaml:"pluralName"`
	Description    string                      `yaml:"description"`
	Attributes     map[string]*AttributeSchema `yaml:"attributes"`
}

type AttributeSchema struct {
	Type string `yaml:"type"`

	Required        bool  `yaml:"required"`
	Max             int64 `yaml:"max"`
	Min             int64 `yaml:"min"`
	MaxLength       int   `yaml:"maxLength"`
	MinLength       int   `yaml:"minLength"`
	Private         bool  `yaml:"private"`
	NotConfigurable bool  `yaml:"notConfigurable"`

	Default     string `yaml:"default"`
	NotNullable bool   `yaml:"notNullable"`
	Unsigned    bool   `yaml:"unsigned"`
	Unique      bool   `yaml:"unique"`
}
