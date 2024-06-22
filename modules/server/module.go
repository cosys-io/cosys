package server

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

func init() {
	svCtor := func(cosys *common.Cosys) common.Server {
		port := cosys.Configs.Server.Port

		return Server{
			Port:  port,
			Cosys: cosys,
		}
	}

	if err := common.RegisterServer("default", svCtor); err != nil {
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
