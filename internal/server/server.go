package server

import (
	"context"
	"net/http"
	"regexp"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/gen/apis"
)

type IServer interface {
	Start() error
}

type Server struct {
	Port  string
	Cosys common.Cosys
}

func NewServer(port string, cosys common.Cosys) *Server {
	return &Server{
		port,
		cosys,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.Handle("/", s)

	err := http.ListenAndServe(s.Port, mux)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, api := range apis.Apis {
		for _, route := range api.Routes {
			matches := route.Regex.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				if route.Method != r.Method {
					continue
				}

				for _, policyName := range route.Policies {
					policy, err := api.Policy(policyName)
					if err != nil {
						// Internal Server Error
						return
					}

					if !policy(s.Cosys, r.Context()) {
						// Rejected
						return
					}
				}

				actionRegex := regexp.MustCompile(`(.+)\.(.+)`)
				controllerAction := actionRegex.FindStringSubmatch(route.Action)
				controllerName := controllerAction[1]
				actionName := controllerAction[2]

				controller, err := api.Controller(controllerName)
				if err != nil {
					// Internal Server Error
					return
				}
				action, err := controller.Action(actionName)
				if err != nil {
					// Internal Server Error
					return
				}

				handler := action(s.Cosys, r.Context())
				for _, middlewareName := range route.Middlewares {
					middleware, err := api.Middleware(middlewareName)
					if err != nil {
						// Internal Server Error
						return
					}
					handler = middleware(s.Cosys, r.Context())(handler)
				}

				ctx := context.WithValue(r.Context(), "query_params", matches[1:])
				handler(w, r.WithContext(ctx))
				return
			}
		}
	}
	http.NotFound(w, r)
}
