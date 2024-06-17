package server

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

func init() {
	svFunc := func(cosys *common.Cosys) common.Server {
		return Server{
			Port:  "3000",
			Cosys: cosys,
		}
	}

	if err := common.RegisterServer("default", svFunc); err != nil {
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
