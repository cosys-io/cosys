package apis

import (
	"context"
	"net/http"

	"github.com/cosys-io/cosys/internal/cosys"
)

type Middleware func(cosys.Cosys, context.Context) func(http.HandlerFunc) http.HandlerFunc
