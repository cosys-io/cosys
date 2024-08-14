package server

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/internal"
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		return cosys.UseServer(internal.NewServer("3000", cosys))
	})
}
