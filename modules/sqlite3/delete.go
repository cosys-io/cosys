package sqlite3

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

	sb.WriteString(model.DBName_())

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

	return sb.String(), nil
}
