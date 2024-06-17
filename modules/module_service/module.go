package module_service

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

func init() {
	msFunc := func(cosys *common.Cosys) common.ModuleService {
		return ModuleService{
			Cosys: cosys,
		}
	}

	if err := common.RegisterModuleService("default", msFunc); err != nil {
		log.Fatal(err)
	}
}

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
