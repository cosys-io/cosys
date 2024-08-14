package internal

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"strconv"
	"strings"
)

func schemaToSQL(schema common.ModelSchema) string {
	var sb strings.Builder

	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(schema.CollectionName())
	sb.WriteString(" ( id INTEGER PRIMARY KEY AUTOINCREMENT, ")

	for index, attr := range schema.Attributes()[1:] {
		sb.WriteString(attr.Name())
		sb.WriteString(" ")
		sb.WriteString(getType(attr.DetailedDataType()))

		var checks []string
		if attr.DetailedDataType() == "Boolean" {
			checks = append(checks, attr.Name()+" IN (0, 1)")
		}
		if attr.Max() != 2147483647 {
			checks = append(checks, attr.Name()+" <= "+strconv.FormatInt(attr.Max(), 10))
		}
		if attr.Min() != -2147483648 {
			checks = append(checks, attr.Name()+" >= "+strconv.FormatInt(attr.Min(), 10))
		}
		if attr.MaxLength() != -1 {
			checks = append(checks, fmt.Sprintf("length(%s) <= %d", attr.Name(), attr.MaxLength()))
		}
		if attr.MinLength() != -1 {
			checks = append(checks, fmt.Sprintf("length(%s) >= %d", attr.Name(), attr.MinLength()))
		}
		if len(checks) > 0 {
			sb.WriteString(" CHECK( ")
			for index, check := range checks {
				sb.WriteString(check)
				if index != len(checks)-1 {
					sb.WriteString(" AND ")
				}
			}
			sb.WriteString(" )")
		}

		if attr.Default() != "" {
			sb.WriteString(" DEFAULT ")
			sb.WriteString(attr.Default())
		}
		if !attr.Nullable() {
			sb.WriteString(" NOT NULL")
		}
		if attr.Unique() {
			sb.WriteString(" UNIQUE")
		}

		if index < len(schema.Attributes())-2 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(" )")

	return sb.String()
}

func getType(attrType string) string {
	switch attrType {
	case "Int":
		return "INTEGER"
	case "Float":
		return "REAL"
	case "Boolean":
		return "INTEGER"
	case "String":
		return "TEXT"
	default:
		return ""
	}
}
