package queryengine

import (
	"strings"

	"github.com/cosys-io/cosys/internal/cosys"
	"github.com/cosys-io/cosys/internal/models"
)

func SQLiteInsertQuery(params *cosys.QEParams, model models.Model) (string, error) {
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")

	sb.WriteString(model.Model_Name())

	insert := params.Columns
	num := len(params.Columns)
	if num == 0 {
		insert = model.Model_All()
		num = len(insert)
	}

	sb.WriteString(" ( ")
	for index, col := range insert {
		insertString := col.DBName()

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
	}

	return sb.String(), nil
}
