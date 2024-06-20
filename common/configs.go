package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
)

type Configs struct {
	Admin    *AdminConfigs
	Database *DatabaseConfigs
	Module   *ModuleConfigs
	Server   *ServerConfigs
}

type AdminConfigs struct{}

type DatabaseConfigs struct {
	Client string `yaml:"client"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

type ModuleConfigs struct {
	Modules []string `yaml:"modules"`
}

type ServerConfigs struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func GetConfigs() (Configs, error) {
	if err := godotenv.Load(); err != nil {
		return Configs{}, err
	}

	adminCfg := &AdminConfigs{}
	databaseCfg := &DatabaseConfigs{}
	moduleCfg := &ModuleConfigs{}
	serverCfg := &ServerConfigs{}

	if err := ParseFile("configs/admin.yaml", adminCfg); err != nil {
		return Configs{}, err
	}
	if err := ParseFile("configs/database.yaml", databaseCfg); err != nil {
		return Configs{}, err
	}
	if err := ParseFile("configs/module.yaml", moduleCfg); err != nil {
		return Configs{}, err
	}
	if err := ParseFile("configs/server.yaml", serverCfg); err != nil {
		return Configs{}, err
	}

	return Configs{
		Admin:    adminCfg,
		Database: databaseCfg,
		Module:   moduleCfg,
		Server:   serverCfg,
	}, nil
}

func ParseFile(path string, obj any) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, obj); err != nil {
		return err
	}

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
