package sqlite3

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/common"
)

func SelectQuery(params *common.DBParams, model common.Model) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}

	var sb strings.Builder

	sb.WriteString("SELECT")

	num := len(params.Select)
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

	sb.WriteString(" FROM ")
	sb.WriteString(model.DBName_())

	num = len(params.Where)
	if num > 0 {
		sb.WriteString(" WHERE ")
		for index, where := range params.Where {
			sb.WriteString("( ")

			whereString, err := StringCondition(where)
			if err != nil {
				return "", err
			}
			sb.WriteString(whereString)

			sb.WriteString(" )")

			if index < num-1 {
				sb.WriteString(" AND ")
			}
		}
	}

	num = len(params.OrderBy)
	if num > 0 {
		sb.WriteString(" ORDER BY ")
		for index, orderBy := range params.OrderBy {

			orderString, err := StringOrder(orderBy)
			if err != nil {
				return "", err
			}
			sb.WriteString(orderString)

			if index < num-1 {
				sb.WriteString(", ")
			}
		}
	}

	sb.WriteString(" LIMIT ")
	sb.WriteString(fmt.Sprint(params.Limit))

	sb.WriteString(" OFFSET ")
	sb.WriteString(fmt.Sprint(params.Offset))

	return sb.String(), nil
}
