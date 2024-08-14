package internal

import (
	"fmt"
	"github.com/cosys-io/cosys/modules/cms/generators"
	"github.com/cosys-io/cosys/modules/cms/schema"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	databaseName string // databaseName is bound to the database flag.
	viewName     string // viewName is bound to the view flag.
	singularName string // singularName is bound to the singular flag.
	pluralName   string // pluralName is bound to the plural flag.
	about        string // about is bound to the about flag.
)

func init() {
	generateCollectionCmd.Flags().StringVarP(&databaseName, "database", "D", "", "name of the sql table for the new content type")
	generateCollectionCmd.Flags().StringVarP(&viewName, "view", "V", "", "name displayed to users for the new content type")
	generateCollectionCmd.Flags().StringVarP(&singularName, "singular", "S", "", "singular name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&pluralName, "plural", "P", "", "plural name of the new content type")
	generateCollectionCmd.Flags().StringVarP(&about, "about", "A", "", "description of the new content type")
	generateCollectionCmd.MarkFlagRequired("singular")
	generateCollectionCmd.MarkFlagRequired("plural")

	generateCmd.AddCommand(generateCollectionCmd)
}

// generateCollectionCmd is the command for generating a collection type.
var generateCollectionCmd = &cobra.Command{
	Use:   "collection [attributes] [flags]",
	Short: "Generate a collection type",
	Long:  "Generate a collection type.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		schema, err := getSchema(databaseName, viewName, singularName, pluralName, about, args)
		if err != nil {
			log.Fatal(err)
		}

		if err := generators.GenerateType(schema); err != nil {
			log.Fatal(err)
		}
	},
}

// getSchema returns the ModelSchema from the given names, description and attribute strings.
func getSchema(databaseName string, viewName string, singularName string, pluralName string,
	about string, attrStrings []string) (*schema.ModelSchema, error) {

	attrs := make([]*schema.AttributeSchema, len(attrStrings)+1)
	attrs[0] = &schema.IdSchema
	attrsSet := make(map[string]bool)

	for index, attrString := range attrStrings {
		attrSchema, err := getAttrSchema(attrString)
		if err != nil {
			return nil, err
		}

		if _, dup := attrsSet[attrSchema.Name()]; dup {
			return nil, fmt.Errorf("duplicate attribute: %s", attrSchema.Name)
		}
		attrsSet[attrSchema.Name()] = true

		attrs[index+1] = attrSchema
	}

	return getModelSchema(databaseName, viewName, singularName, pluralName, about, attrs), nil
}

// getType returns the simple and detailed type from the given attribute type string.
func getType(attrType string) (string, string, error) {
	var attrSimpleType string
	var attrDetailedType string

	switch attrType {
	case "int":
		attrSimpleType = "Number"
		attrDetailedType = "Int"
	case "float":
		attrSimpleType = "Number"
		attrDetailedType = "Float"
	case "bool":
		attrSimpleType = "Boolean"
		attrDetailedType = "Boolean"
	case "string":
		attrSimpleType = "String"
		attrDetailedType = "String"
	case "date":
		attrSimpleType = "Date"
		attrDetailedType = "Date"
	case "datetime":
		attrSimpleType = "DateTime"
		attrDetailedType = "DateTime"
	case "timestamp":
		attrSimpleType = "Timestamp"
		attrDetailedType = "Timestamp"
	default:
		return "", "", fmt.Errorf("invalid type: %s", attrType)
	}

	return attrSimpleType, attrDetailedType, nil
}

// getModelSchema returns the ModelSchema from the given names, description and attribute schemas.
func getModelSchema(databaseName, viewName, singularName, pluralName, about string, attrs []*schema.AttributeSchema) *schema.ModelSchema {
	var collectionName string
	if databaseName != "" {
		collectionName = strcase.ToLowerCamel(databaseName)
	} else {
		collectionName = strcase.ToLowerCamel(pluralName)
	}

	var displayName string
	if viewName != "" {
		displayName = strcase.ToDelimited(viewName, ' ')
	} else {
		displayName = strcase.ToDelimited(pluralName, ' ')
	}

	return schema.NewModelSchema(collectionName, displayName, singularName, pluralName, about, attrs...)
}

// getAttrSchema returns the AttributeSchema from the given attribute string.
func getAttrSchema(attrString string) (*schema.AttributeSchema, error) {
	split := strings.Split(attrString, ":")
	if len(split) < 2 {
		return nil, fmt.Errorf("invalid attribute format: %s", attrString)
	}
	attrName := strcase.ToLowerCamel(split[0])
	attrType := split[1]

	attrSimpleType, attrDetailedType, err := getType(attrType)
	if err != nil {
		return nil, err
	}

	attrSchema := schema.NewAttrSchema(attrName, attrSimpleType, attrDetailedType)

	for _, optionString := range split[2:] {
		option, err := getOption(optionString)
		if err != nil {
			return nil, err
		}
		option(attrSchema)
	}

	return attrSchema, nil
}

// getOption returns the attribute configuration from the given attribute configuration string.
func getOption(option string) (schema.AttrOption, error) {
	switch {
	case option == "notshown":
		return schema.HideInTable, nil
	case option == "required":
		return schema.Required, nil
	case regexp.MustCompile(`^max=([-0-9]+)$`).MatchString(option):
		matches := regexp.MustCompile(`^max=([-0-9]+)$`).FindStringSubmatch(option)
		val, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return schema.Max(val), nil
	case regexp.MustCompile(`^min=([-0-9]+)$`).MatchString(option):
		matches := regexp.MustCompile(`^min=([-0-9]+)$`).FindStringSubmatch(option)
		val, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return schema.Min(val), nil
	case regexp.MustCompile(`^maxlength=(\d+)$`).MatchString(option):
		matches := regexp.MustCompile(`^maxlength=(\d+)$`).FindStringSubmatch(option)
		val, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}
		return schema.MaxLength(val), nil
	case regexp.MustCompile(`^minlength=(\d+)$`).MatchString(option):
		matches := regexp.MustCompile(`^minlength=(\d+)$`).FindStringSubmatch(option)
		val, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}
		return schema.MinLength(val), nil
	case option == "private":
		return schema.Private, nil
	case option == "noteditable":
		return schema.NotEditable, nil
	case regexp.MustCompile(`^default=(.+)=$`).MatchString(option):
		matches := regexp.MustCompile(`^default=(.+)$`).FindStringSubmatch(option)
		return schema.Default(matches[1]), nil
	case option == "notnullable":
		return schema.NotNullable, nil
	case option == "unique":
		return schema.Unique, nil
	default:
		return nil, fmt.Errorf("invalid option: %s", option)
	}
}
