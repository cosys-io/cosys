package queryengine

import (
	"fmt"
	"strings"

	"github.com/cosys-io/cosys/internal/cosys"
	"github.com/cosys-io/cosys/internal/models"
)

func SQLiteDeleteQuery(params *cosys.QEParams, model models.Model) (string, error) {
	var sb strings.Builder

	sb.WriteString("DELETE FROM ")

	sb.WriteString(model.Model_Name())

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
