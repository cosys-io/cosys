package internal

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/common"
)

// deleteQuery returns an insert sql query from the given params.
func insertQuery(params *common.DBParams, model common.Model) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}

	var sb strings.Builder

	sb.WriteString("INSERT INTO ")

	sb.WriteString(model.PluralSnakeName_())

	insert := params.Columns
	num := len(params.Columns)
	if num == 0 {
		insert = model.Attributes_()[1:]
		num = len(insert)
	}

	sb.WriteString(" ( ")
	for index, col := range insert {
		insertString := col.SnakeName()

		sb.WriteString(insertString)
		if index < num-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" )")

	sb.WriteString(" VALUES ( ")
	for index := range num {
		sb.WriteString("?")
		if index < num-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" )")

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
