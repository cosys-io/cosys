package content_builder

import "github.com/cosys-io/cosys/common"

var Module = &common.Module{
	Routes: nil,
	Controllers: map[string]common.Controller{
		"admin": Controller,
	},
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: nil,
	OnDestroy:  nil,
}
