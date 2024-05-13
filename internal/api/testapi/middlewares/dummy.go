package middlewares

import (
	"context"
	"net/http"

	"github.com/cosys-io/cosys/internal/cosys"
)

func DummyMiddleware(cs cosys.Cosys, ctx context.Context) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
