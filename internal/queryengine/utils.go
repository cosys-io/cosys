package queryengine

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/cosys-io/cosys/internal/cosys"
	"github.com/cosys-io/cosys/internal/models"
)

func Extract(data models.Entity, params *cosys.QEParams, model models.Model) ([]any, error) {
	attributes := []any{}

	columns := params.Columns
	if len(params.Columns) == 0 {
		columns = model.Model_All()
	}

	var dataValue reflect.Value = reflect.ValueOf(data)
	if reflect.TypeOf(data).Kind() == reflect.Pointer {
		dataValue = reflect.Indirect(dataValue)
	}
	for _, col := range columns {
		name := col.StructName()

		attributeValue := dataValue.FieldByName(name)
		if attributeValue == (reflect.Value{}) {
			return nil, fmt.Errorf("attribute %s not found", name)
		}

		attributes = append(attributes, attributeValue.Interface())
	}

	return attributes, nil
}

func Scan(rows *sql.Rows, params *cosys.QEParams, model models.Model) (models.Entity, error) {
	entity := model.Model_New()

	selects := params.Selects
	if len(params.Selects) == 0 {
		selects = model.Model_All()
	}

	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity).Convert(entityType)
	if entityType.Kind() == reflect.Pointer {
		entityValue = reflect.Indirect(entityValue)
	}

	numCols := len(selects)
	columns := make([]any, numCols)
	for index, attribute := range selects {
		field := entityValue.FieldByName(attribute.StructName())
		columns[index] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func SelectQuery(dialect string, params *cosys.QEParams, model models.Model) (string, error) {
	switch dialect {
	case "sqlite3":
		return SQLiteSelectQuery(params, model)
	default:
		return "", fmt.Errorf("dialect not supported: %s", dialect)
	}
}

func InsertQuery(dialect string, params *cosys.QEParams, model models.Model) (string, error) {
	switch dialect {
	case "sqlite3":
		return SQLiteInsertQuery(params, model)
	default:
		return "", fmt.Errorf("dialect not supported: %s", dialect)
	}
}

func UpdateQuery(dialect string, params *cosys.QEParams, model models.Model) (string, error) {
	switch dialect {
	case "sqlite3":
		return SQLiteUpdateQuery(params, model)
	default:
		return "", fmt.Errorf("dialect not supported: %s", dialect)
	}
}

func DeleteQuery(dialect string, params *cosys.QEParams, model models.Model) (string, error) {
	switch dialect {
	case "sqlite3":
		return SQLiteDeleteQuery(params, model)
	default:
		return "", fmt.Errorf("dialect not supported: %s", dialect)
	}
}

func StringCondition(where models.Condition) (string, error) {
	switch where := where.(type) {
	case *models.NestedCondition:
		return StringNested(where)
	case *models.ExpressionCondition:
		return StringExpressions(where)
	case *models.BoolAttribute:
		return where.DBName(), nil
	default:
		return "", fmt.Errorf("invalid where condition")
	}
}

func StringNested(where *models.NestedCondition) (string, error) {
	var err error

	var left string
	if where.Left == nil {
		return "", fmt.Errorf("left operand not found")
	} else {
		left, err = StringCondition(where.Left)
		if err != nil {
			return "", err
		}
	}

	var right string
	if where.Right == nil {
		right = ""
	} else {
		right, err = StringCondition(where.Right)
		if err != nil {
			return "", err
		}
	}

	switch where.Op {
	case models.NOT:
		return fmt.Sprintf("NOT ( %s )", left), nil
	case models.AND, models.OR:
		if right == "" {
			return "", fmt.Errorf("right operand not found")
		}

		return fmt.Sprintf("( %s ) %s ( %s )", left, where.Op, right), err
	default:
		return "", fmt.Errorf("illegal operation: %s", where.Op)
	}
}

func StringExpressions(where *models.ExpressionCondition) (string, error) {
	var left string
	if where.Left == nil {
		return "", fmt.Errorf("right operand not found")
	} else {
		left = where.Left.DBName()
	}

	switch where.Op {
	case models.NONE:
		return left, nil
	case models.EQ, models.NEQ:
		switch r := where.Right.(type) {
		case string:
			return fmt.Sprintf(`%s %s "%s"`, left, where.Op, r), nil
		case int:
			right := strconv.Itoa(r)

			return fmt.Sprintf("%s %s %s", left, where.Op, right), nil
		case bool:
			right := strconv.FormatBool(r)

			return fmt.Sprintf("%s %s %s", left, where.Op, right), nil
		default:
			return "", fmt.Errorf("illegal right operand: %s", where.Right)
		}
	case models.LT, models.GT, models.LTE, models.GTE:
		if _, ok := where.Right.(int); !ok {
			return "", fmt.Errorf("illegal right operand: %s", where.Right)
		}
		right := strconv.Itoa(where.Right.(int))

		return fmt.Sprintf("%s %s %s", left, where.Op, right), nil
	case models.NULL, models.NOTNULL:
		return fmt.Sprintf("%s %s", left, where.Op), nil
	default:
		return "", fmt.Errorf("illegal operation: %s", where.Op)
	}
}

func StringOrder(orderBy *models.Order) (string, error) {
	if orderBy.Order == models.ASC || orderBy.Order == models.DESC {
		return fmt.Sprintf("%s %s", orderBy.Attribute.DBName(), orderBy.Order), nil
	}
	return "", fmt.Errorf("illegal operation for order condition: %s", orderBy.Order)
}
