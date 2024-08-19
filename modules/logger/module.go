package logger

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/logger/internal"
)

// init registers the module to register the Logger core service.
func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		return cosys.UseLogger(internal.Logger{})
	})
}
