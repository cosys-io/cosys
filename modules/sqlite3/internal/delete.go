package internal

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/common"
)

func DeleteQuery(params *common.DBParams, model common.Model) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}

	var sb strings.Builder

	sb.WriteString("DELETE FROM ")

	sb.WriteString(model.PluralSnakeName_())

	num := len(params.Where)
	if num > 0 {
		sb.WriteString(" WHERE")
		for index, where := range params.Where {
			sb.WriteString(" ")

			whereString, err := StringCondition(where)
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
