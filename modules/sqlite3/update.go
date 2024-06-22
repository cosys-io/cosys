package sqlite3

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/common"
)

func UpdateQuery(params *common.DBParams, model common.Model) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}

	var sb strings.Builder

	sb.WriteString("UPDATE ")

	sb.WriteString(model.Name_())

	sb.WriteString(" SET ")

	update := params.Columns
	num := len(params.Columns)
	if num == 0 {
		update = model.All_()[1:]
		num = len(update)
	}

	for i := range num {
		sb.WriteString(update[i].Name())
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

	return sb.String(), nil
}
