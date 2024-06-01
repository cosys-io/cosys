package queryengine

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/models"
)

func SQLiteSelectQuery(params *common.QEParams, model models.Model) (string, error) {
	var sb strings.Builder

	sb.WriteString("SELECT")

	num := len(params.Selects)
	if num == 0 {
		sb.WriteString(" *")
	} else {
		for index, col := range params.Selects {
			sb.WriteString(" ")

			colString := col.Name()

			sb.WriteString(colString)

			if index < num-1 {
				sb.WriteString(",")
			}
		}
	}

	sb.WriteString(" FROM ")
	sb.WriteString(model.Name_())

	num = len(params.Wheres)
	if num > 0 {
		sb.WriteString(" WHERE ")
		for index, where := range params.Wheres {
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

	num = len(params.OrderBys)
	if num > 0 {
		sb.WriteString(" ORDER BY ")
		for index, orderBy := range params.OrderBys {

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
	sb.WriteString(fmt.Sprint(params.LimitVal))

	sb.WriteString(" OFFSET ")
	sb.WriteString(fmt.Sprint(params.OffsetVal))

	return sb.String(), nil
}
