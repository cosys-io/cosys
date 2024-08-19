package internal

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/response"
	"net/http"
)

// Server is an implementation of the Server core service using the native net/http package.
type Server struct {
	port  string
	mux   *http.ServeMux
	cosys *common.Cosys
}

// NewServer returns a new Server.
func NewServer(port string, cosys *common.Cosys) *Server {
	return &Server{
		port:  port,
		mux:   new(http.ServeMux),
		cosys: cosys,
	}
}

// resolveEndpoints creates the mux from the registered routes, controllers, middlewares and policies.
func (s Server) resolveEndpoints() error {
	mux := http.NewServeMux()

	routes := s.cosys.Routes()

	for _, route := range routes {
		handleFunc, err := route.Action(s.cosys)
		if err != nil {
			return err
		}

		for i := len(route.Middlewares) - 1; i > 0; i-- {
			middleware, err := route.Middlewares[i](s.cosys)
			if err != nil {
				return err
			}

			handleFunc = middleware(handleFunc)
		}

		for i := len(route.Policies) - 1; i > 0; i-- {
			policy, err := route.Policies[i](s.cosys)
			if err != nil {
				return err
			}

			policyMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if !policy(r) {
						response.RespondError(w, "Forbidden", http.StatusForbidden)
					}

					next.ServeHTTP(w, r)
				}
			}

			handleFunc = policyMiddleware(handleFunc)
		}

		mux.HandleFunc(route.Method+" "+route.Path, handleFunc)
	}

	*s.mux = *mux
	return nil
}

// Start resolved the server endpoints and starts the server.
func (s Server) Start() error {
	if err := s.resolveEndpoints(); err != nil {
		return err
	}
	if err := http.ListenAndServe(":"+s.port, s.mux); err != nil {
		return err
	}

	return nil
}
