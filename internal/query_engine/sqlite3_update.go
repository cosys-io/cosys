package query_engine

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/models"
)

func SQLiteUpdateQuery(params *common.QEParams, model models.Model) (string, error) {
	var sb strings.Builder

	sb.WriteString("UPDATE ")

	sb.WriteString(model.Name_())

	sb.WriteString(" SET ")

	update := params.Columns
	num := len(params.Columns)
	if num == 0 {
		update = model.All_()
		num = len(update)
	}

	for i := range num {
		sb.WriteString(update[i].Name())
		sb.WriteString(" = ?")
		if i < num-1 {
			sb.WriteString(", ")
		}
	}

	num = len(params.Wheres)
	if num > 0 {
		sb.WriteString(" WHERE")
		for index, where := range params.Wheres {
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
