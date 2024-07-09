package common

import (
	"github.com/joho/godotenv"
	"path/filepath"
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

func GetConfigs(path string) (Configs, error) {
	if err := godotenv.Load(); err != nil {
		return Configs{}, err
	}

	adminCfg := AdminConfigs{}
	databaseCfg := DatabaseConfigs{}
	moduleCfg := ModuleConfigs{}
	serverCfg := ServerConfigs{}

	if err := ParseFile(filepath.Join(path, "admin.yaml"), &adminCfg, true); err != nil {
		return Configs{}, err
	}
	if err := ParseFile(filepath.Join(path, "database.yaml"), &databaseCfg, true); err != nil {
		return Configs{}, err
	}
	if err := ParseFile(filepath.Join(path, "module.yaml"), &moduleCfg, true); err != nil {
		return Configs{}, err
	}
	if err := ParseFile(filepath.Join(path, "server.yaml"), &serverCfg, true); err != nil {
		return Configs{}, err
	}

	return Configs{
		Admin:    &adminCfg,
		Database: &databaseCfg,
		Module:   &moduleCfg,
		Server:   &serverCfg,
	}, nil
}

func NewConfigs() Configs {
	return Configs{
		Admin: &AdminConfigs{},
		Database: &DatabaseConfigs{
			Client: "sqlite3",
			Host:   "localhost",
			Port:   "4000",
			Name:   "cosys",
			User:   "cosys",
			Pass:   "cosys",
		},
		Module: &ModuleConfigs{
			Modules: []string{},
		},
		Server: &ServerConfigs{
			Host: "localhost",
			Port: "3000",
		},
	}
}
