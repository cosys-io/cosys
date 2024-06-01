package middlewares

import (
	"github.com/cosys-io/cosys/internal/apis"
)

var Middlewares = map[string]apis.Middleware{
	"dummy": DummyMiddleware,
}
