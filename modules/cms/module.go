package cms

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/cms/internal"
	"github.com/spf13/cobra"
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		return cosys.AddCommands(func(*common.Cosys) *cobra.Command { return internal.RootCmd })
	})
}
