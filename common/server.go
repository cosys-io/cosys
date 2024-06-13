package common

import (
	"context"
	"net/http"
	"regexp"
)

type Server struct {
	Port  string
	Cosys *Cosys
}

func NewServer(port string, cosys *Cosys) *Server {
	return &Server{
		port,
		cosys,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, module := range s.Cosys.Modules {
			for _, route := range module.Routes {
				matches := route.Regex.FindStringSubmatch(r.URL.Path)
				if len(matches) > 0 {
					if r.Method != route.Method {
						continue
					}
					ctx := context.WithValue(r.Context(), "query_params", matches[1:])

					for _, policyName := range route.Policies {
						policy, ok := module.Policies[policyName]
						if !ok {
							// Internal Server Error
							return
						}
						if !policy(*s.Cosys, ctx) {
							// Forbidden
							return
						}
					}

					uidRegex := regexp.MustCompile(`(.+)\.(.+)`)
					actionUid := route.Action
					uidMatches := uidRegex.FindStringSubmatch(actionUid)
					if len(uidMatches) != 3 {
						// Internal Server Error
						return
					}
					controllerName := uidMatches[1]
					actionName := uidMatches[2]
					controller, ok := module.Controllers[controllerName]
					if !ok {
						// Internal Server Error
						return
					}
					actionFunc, ok := controller.Actions[actionName]
					if !ok {
						// Internal Server Error
						return
					}
					action := actionFunc(*s.Cosys, ctx)

					for _, middlewareName := range route.Middlewares {
						middlewareFunc, ok := module.Middlewares[middlewareName]
						if !ok {
							// Internal Server Error
							return
						}
						middleware := middlewareFunc(*s.Cosys, ctx)

						action = middleware(action)
					}

					action(w, r.WithContext(ctx))
					return
				}
			}
		}
		http.NotFound(w, r)
	})

	if err := http.ListenAndServe(":"+s.Port, mux); err != nil {
		return err
	}

	return nil
}
