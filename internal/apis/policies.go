package apis

import (
	"context"

	"github.com/cosys-io/cosys/internal/cosys"
)

type Policy func(cosys.Cosys, context.Context) bool
