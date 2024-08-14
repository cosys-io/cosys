package server

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/internal"
)

// init registers the module to register the Server core service.
func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		return cosys.UseServer(internal.NewServer("3000", cosys))
	})
}
