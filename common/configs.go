package common

type Configs struct {
	Admin    *AdminConfigs
	Database *DatabaseConfigs
	Server   *ServerConfigs
}

type AdminConfigs struct{}

type AdminConfigsFunc func(*Env) *AdminConfigs

type DatabaseConfigs struct{}

type DatabaseConfigsFunc func(*Env) *DatabaseConfigs

type ServerConfigs struct{}

type ServerConfigsFunc func(*Env) *ServerConfigs

type Env struct{}
