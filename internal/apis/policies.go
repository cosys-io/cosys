package apis

import (
	"context"

	"github.com/cosys-io/cosys/internal/common"
)

type Policy func(common.Cosys, context.Context) bool
