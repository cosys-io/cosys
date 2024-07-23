package logger

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/logger/internal"
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		return cosys.UseLogger(internal.Logger{})
	})
}
