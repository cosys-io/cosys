package internal

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/common"
)

// updateQuery returns an update sql query from the given params.
func updateQuery(params *common.DBParams, model common.Model) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}

	var sb strings.Builder

	sb.WriteString("UPDATE ")

	sb.WriteString(model.PluralSnakeName_())

	sb.WriteString(" SET ")

	update := params.Columns
	num := len(params.Columns)
	if num == 0 {
		update = model.Attributes_()[1:]
		num = len(update)
	}

	for i := range num {
		sb.WriteString(update[i].SnakeName())
		sb.WriteString(" = ?")
		if i < num-1 {
			sb.WriteString(", ")
		}
	}

	num = len(params.Where)
	if num > 0 {
		sb.WriteString(" WHERE")
		for index, where := range params.Where {
			sb.WriteString(" ")

			whereString, err := stringCondition(where)
			if err != nil {
				return "", err
			}

			sb.WriteString(whereString)
			if index < num-1 {
				sb.WriteString(", ")
			}
		}
	} else {
		return "", fmt.Errorf("where condition not found")
	}

	sb.WriteString(" RETURNING")

	num = len(params.Select)
	if num == 0 {
		sb.WriteString(" *")
	} else {
		for index, col := range params.Select {
			sb.WriteString(" ")

			colString := col.SnakeName()

			sb.WriteString(colString)

			if index < num-1 {
				sb.WriteString(",")
			}
		}
	}

	return sb.String(), nil
}
