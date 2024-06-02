package query_engine

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/models"
)

func SQLiteDeleteQuery(params *common.QEParams, model models.Model) (string, error) {
	var sb strings.Builder

	sb.WriteString("DELETE FROM ")

	sb.WriteString(model.Name_())

	num := len(params.Wheres)
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
