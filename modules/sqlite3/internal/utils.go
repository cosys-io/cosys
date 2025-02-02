package internal

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/cosys-io/cosys/common"
)

// extract returns a slice of the values of an entity's fields.
func extract(data common.Entity, params *common.DBParams, model common.Model) ([]any, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	var attrs []any

	columns := params.Columns
	if len(params.Columns) == 0 {
		columns = model.Attributes_()[1:]
	}

	dataValue := reflect.ValueOf(data)
	if reflect.TypeOf(data).Kind() == reflect.Pointer {
		dataValue = reflect.Indirect(dataValue)
	}
	if dataValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}
	for _, col := range columns {
		attrName := col.PascalName()

		attributeValue := dataValue.FieldByName(attrName)
		if attributeValue == (reflect.Value{}) {
			return nil, fmt.Errorf("attribute not found: %s", attrName)
		}

		attrs = append(attrs, attributeValue.Interface())
	}

	return attrs, nil
}

// scan scans the values from a sql row into an entity struct.
func scan(rows *sql.Rows, params *common.DBParams, model common.Model) (common.Entity, error) {
	if rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}

	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	entity := model.New_()

	selects := params.Select
	if len(params.Select) == 0 {
		selects = model.Attributes_()
	}

	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity).Convert(entityType)
	if entityType.Kind() == reflect.Pointer {
		entityValue = reflect.Indirect(entityValue)
	}
	if entityValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity is not a struct")
	}

	numCols := len(selects)
	columns := make([]any, numCols)
	for index, attribute := range selects {
		field := entityValue.FieldByName(attribute.PascalName())
		columns[index] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// stringCondition returns the sql condition for a where condition.
func stringCondition(where common.Condition) (string, error) {
	if where == nil {
		return "", fmt.Errorf("where is nil")
	}

	switch where := where.(type) {
	case *common.NestedCondition:
		return stringNested(where)
	case *common.ExpressionCondition:
		return stringExpressions(where)
	case *common.BoolAttribute:
		return where.SnakeName(), nil
	default:
		return "", fmt.Errorf("invalid where condition")
	}
}

// stringNested returns the sql condition for a nested where condition.
func stringNested(where *common.NestedCondition) (string, error) {
	if where == nil {
		return "", fmt.Errorf("where is nil")
	}

	var err error

	var left string
	if where.Left == nil {
		return "", fmt.Errorf("left operand not found")
	} else {
		left, err = stringCondition(where.Left)
		if err != nil {
			return "", err
		}
	}

	var right string
	if where.Right == nil {
		right = ""
	} else {
		right, err = stringCondition(where.Right)
		if err != nil {
			return "", err
		}
	}

	switch where.Op {
	case common.Not:
		return fmt.Sprintf("NOT ( %s )", left), nil
	case common.And, common.Or:
		if right == "" {
			return "", fmt.Errorf("right operand not found")
		}

		return fmt.Sprintf("( %s ) %s ( %s )", left, where.Op, right), err
	default:
		return "", fmt.Errorf("illegal operation: %s", where.Op)
	}
}

// stringExpressions returns the sql condition for an expression where condition.
func stringExpressions(where *common.ExpressionCondition) (string, error) {
	if where == nil {
		return "", fmt.Errorf("where is nil")
	}

	var left string
	if where.Left == nil {
		return "", fmt.Errorf("right operand not found")
	} else {
		left = where.Left.SnakeName()
	}

	switch where.Op {
	case common.None:
		return left, nil
	case common.Eq, common.Neq:
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
	case common.Lt, common.Gt, common.Lte, common.Gte:
		if _, ok := where.Right.(int); !ok {
			return "", fmt.Errorf("illegal right operand: %s", where.Right)
		}
		right := strconv.Itoa(where.Right.(int))

		return fmt.Sprintf("%s %s %s", left, where.Op, right), nil
	case common.Null, common.NotNull:
		return fmt.Sprintf("%s %s", left, where.Op), nil
	default:
		return "", fmt.Errorf("illegal operation: %s", where.Op)
	}
}

// stringOrder returns the sql condition for an order condition.
func stringOrder(orderBy *common.Order) (string, error) {
	if orderBy == nil {
		return "", fmt.Errorf("order is nil")
	}

	if orderBy.Order == common.Asc || orderBy.Order == common.Desc {
		return fmt.Sprintf("%s %s", orderBy.Attribute.SnakeName(), orderBy.Order), nil
	}
	return "", fmt.Errorf("illegal operation for order condition: %s", orderBy.Order)
}
