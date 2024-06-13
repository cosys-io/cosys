package cmd

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	displayName  string
	singularName string
	pluralName   string
	description  string
)

func init() {
	generateCollectionCmd.Flags().StringVarP(&displayName, "display", "N", "", "display name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&singularName, "singular", "S", "", "singular name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&pluralName, "plural", "P", "", "plural name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&description, "description", "D", "", "description of the new content type")

	generateCmd.AddCommand(generateCollectionCmd)
}

var generateCollectionCmd = &cobra.Command{
	Use:   "collection content_type_name ",
	Short: "Generate a collection type",
	Long:  "Generate a collection type.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collectionName := args[0]

		attributes := args[1:]

		schema, err := schemaFromArgs(collectionName, displayName, singularName, pluralName, description, attributes)
		if err != nil {
			log.Fatal(err)
		}

		if err := generateType(collectionName, schema); err != nil {
			log.Fatal(err)
		}
	},
}

func generateType(typeName string, schema *common.ModelSchema) error {
	typeDir := filepath.Join("modules/api/content_types", typeName)

	if err := generateDir(typeDir, genHeadOnly); err != nil {
		return err
	}

	ctx := CtxFromSchema(schema)

	if err := generateFile(filepath.Join(typeDir, "schema.yaml"), SchemaYamlTmpl, schema); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "schema.go"), SchemaGoTmpl, schema); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "lifecycle.go"), LifecycleTmpl, ctx); err != nil {
		return err
	}

	if err := generateFile(filepath.Join(typeDir, "model.go"), ModelTmpl, ctx); err != nil {
		return err
	}

	return nil
}

func schemaFromArgs(collection string, display string, singular string, plural string, description string, attrStrings []string) (*common.ModelSchema, error) {
	modelSchema := &common.ModelSchema{
		ModelType:      "collectionType",
		CollectionName: collection,
		DisplayName:    display,
		SingularName:   singular,
		PluralName:     plural,
		Description:    description,
		Attributes:     map[string]*common.AttributeSchema{},
	}

	for _, attrString := range attrStrings {
		split := strings.Split(attrString, ":")
		if len(split) < 2 {
			return nil, fmt.Errorf("invalid attribute format: %s", attrString)
		}
		attrName := split[0]
		attrType := split[1]

		attrSchema := &common.AttributeSchema{
			Type:            attrType,
			Required:        false,
			Max:             0,
			Min:             0,
			MaxLength:       0,
			MinLength:       0,
			Private:         false,
			NotConfigurable: false,
			Default:         "",
			NotNullable:     false,
			Unsigned:        false,
			Unique:          false,
		}

		for _, option := range split[2:] {
			switch {
			case option == "required":
				attrSchema.Required = true
			case regexp.MustCompile(`^max\(([-.0-9]+)\)$`).MatchString(option):
				matches := regexp.MustCompile(`^max\(([-.0-9]+)\)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, err
				}
				attrSchema.Max = val
			case regexp.MustCompile(`^min\(([-.0-9]+)\)$`).MatchString(option):
				matches := regexp.MustCompile(`^min\(([-.0-9]+)\)$`).FindStringSubmatch(option)
				val, err := strconv.ParseInt(matches[1], 10, 64)
				if err != nil {
					return nil, err
				}
				attrSchema.Min = val
			case regexp.MustCompile(`^maxlength\((\d+)\)$`).MatchString(option):
				matches := regexp.MustCompile(`^maxlength\((\d+)\)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				attrSchema.MaxLength = val
			case regexp.MustCompile(`^minlength\((\d+)\)$`).MatchString(option):
				matches := regexp.MustCompile(`^minlength\((\d+)\)$`).FindStringSubmatch(option)
				val, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				attrSchema.MinLength = val
			case option == "private":
				attrSchema.Private = true
			case option == "notconfigurable":
				attrSchema.NotConfigurable = true
			case regexp.MustCompile(`^default\((.+)\)$`).MatchString(option):
				matches := regexp.MustCompile(`^default\((.+)\)$`).FindStringSubmatch(option)
				attrSchema.Default = matches[1]
			case option == "notnullable":
				attrSchema.NotNullable = true
			case option == "unsigned":
				attrSchema.Unsigned = true
			case option == "unique":
				attrSchema.Unique = true
			}
		}

		if _, ok := modelSchema.Attributes[attrName]; ok {
			return nil, fmt.Errorf("duplicate attribute: %s", attrName)
		}
		modelSchema.Attributes[attrName] = attrSchema
	}

	return modelSchema, nil
}

type ModelCtx struct {
	CollectionName string
	SingularName   string
	PluralName     string
	Attributes     []*AttributeCtx
}

type AttributeCtx struct {
	TypeLower  string
	TypeUpper  string
	NameCamel  string
	NamePascal string
}

func CtxFromSchema(schema *common.ModelSchema) *ModelCtx {
	caser := cases.Title(language.English)

	modelCtx := &ModelCtx{
		CollectionName: schema.CollectionName,
		SingularName:   caser.String(schema.SingularName),
		PluralName:     caser.String(schema.PluralName),
		Attributes: []*AttributeCtx{
			{
				TypeLower:  "int",
				TypeUpper:  "Int",
				NameCamel:  "id",
				NamePascal: "Id",
			},
		},
	}

	for name, attr := range schema.Attributes {
		attrCtx := &AttributeCtx{
			NameCamel:  name,
			NamePascal: caser.String(name),
		}

		switch attr.Type {
		case "integer":
			attrCtx.TypeLower = "int"
			attrCtx.TypeUpper = "Int"
		case "string":
			attrCtx.TypeLower = "string"
			attrCtx.TypeUpper = "String"
		case "boolean":
			attrCtx.TypeLower = "bool"
			attrCtx.TypeUpper = "Bool"
		}

		modelCtx.Attributes = append(modelCtx.Attributes, attrCtx)
	}

	return modelCtx
}
