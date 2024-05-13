package policies

import (
	"context"

	"github.com/cosys-io/cosys/internal/cosys"
)

func DummyPolicy(cs cosys.Cosys, ctx context.Context) bool {
	return true
}
