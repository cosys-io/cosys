package sqlite3

import "github.com/cosys-io/cosys/common"

var Module = &common.Module{
	Routes:      nil,
	Controllers: nil,
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: Register,
	OnDestroy:  nil,
}
