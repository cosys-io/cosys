package internal

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/response"
	"net/http"
)

type Server struct {
	port  string
	mux   *http.ServeMux
	cosys *common.Cosys
}

func NewServer(port string, cosys *common.Cosys) *Server {
	return &Server{
		port:  port,
		mux:   new(http.ServeMux),
		cosys: cosys,
	}
}

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

func (s Server) Start() error {
	if err := s.resolveEndpoints(); err != nil {
		return err
	}
	if err := http.ListenAndServe(":"+s.port, s.mux); err != nil {
		return err
	}

	return nil
}
