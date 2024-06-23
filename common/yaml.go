package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

func ParseFile(path string, obj any, getEnv bool) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, obj); err != nil {
		return err
	}

	if getEnv {
		if err = GetEnv(obj); err != nil {
			return err
		}
	}

	return nil
}

func GetEnv(obj any) error {
	objValue := reflect.ValueOf(obj)
	if reflect.TypeOf(obj).Kind() == reflect.Pointer {
		objValue = reflect.Indirect(objValue)
	}
	if objValue.Kind() != reflect.Struct {
		return fmt.Errorf("object is not a struct")
	}

	numFields := objValue.NumField()
	for i := 0; i < numFields; i++ {
		fieldValue := objValue.Field(i)
		if !fieldValue.IsValid() {
			return fmt.Errorf("object has no field %s", objValue.Type().Field(i).Name)
		}
		field := fieldValue.Interface()
		if reflect.TypeOf(field).String() != "string" {
			continue
		}
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s cannot be set", objValue.Type().Field(i).Name)
		}
		fieldString := field.(string)
		if len(fieldString) < 5 {
			continue
		}
		if fieldString[:4] == "ENV." {
			fieldEnv := os.Getenv(fieldString[4:])
			fieldValue.SetString(fieldEnv)
		}
	}

	return nil
}
