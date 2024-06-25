package logger

import (
	"github.com/cosys-io/cosys/common"
	_ "github.com/cosys-io/cosys/modules/logger/common"
)

var Module = &common.Module{
	Routes:      nil,
	Controllers: nil,
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: nil,
	OnDestroy:  nil,
}
