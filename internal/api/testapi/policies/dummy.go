package policies

import (
	"context"

	"github.com/cosys-io/cosys/internal/common"
)

func DummyPolicy(cs common.Cosys, ctx context.Context) bool {
	return true
}
