package common

type ModelSchema struct {
	ModelType      string             `yaml:"modelType" json:"modelType"`
	CollectionName string             `yaml:"collectionName" json:"collectionName"`
	DisplayName    string             `yaml:"displayName" json:"displayName"`
	SingularName   string             `yaml:"singularName" json:"singularName"`
	PluralName     string             `yaml:"pluralName" json:"pluralName"`
	Description    string             `yaml:"description" json:"description"`
	Attributes     []*AttributeSchema `yaml:"attributes" json:"attributes"`
}

type AttributeSchema struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`

	Required        bool  `yaml:"required" json:"required"`
	Max             int64 `yaml:"max" json:"max"`
	Min             int64 `yaml:"min" json:"min"`
	MaxLength       int   `yaml:"maxLength" json:"maxLength"`
	MinLength       int   `yaml:"minLength" json:"minLength"`
	Private         bool  `yaml:"private" json:"private"`
	NotConfigurable bool  `yaml:"notConfigurable" json:"notConfigurable"`

	Default     string `yaml:"default" json:"default"`
	NotNullable bool   `yaml:"notNullable" json:"notNullable"`
	Unsigned    bool   `yaml:"unsigned" json:"unsigned"`
	Unique      bool   `yaml:"unique" json:"unique"`
}
