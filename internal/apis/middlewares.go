package apis

import (
	"context"
	"net/http"

	"github.com/cosys-io/cosys/internal/common"
)

type Middleware func(common.Cosys, context.Context) func(http.HandlerFunc) http.HandlerFunc
