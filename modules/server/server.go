package server

import (
	"bytes"
	"context"
	"github.com/cosys-io/cosys/common"
	"net/http"
	"regexp"
)

type Server struct {
	Port  string
	Cosys *common.Cosys
}

func (s Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, module := range s.Cosys.Modules {
			for _, route := range module.Routes {
				matches := route.Regex.FindStringSubmatch(r.URL.Path)
				if len(matches) > 0 {
					if r.Method != route.Method {
						continue
					}

					var log bytes.Buffer
					w = common.ResponseWriter{
						Writer: w,
						Log:    log,
					}
					ctx := context.WithValue(r.Context(), common.ResponseKey, &log)
					ctx = context.WithValue(ctx, common.StateKey, map[string]any{
						"query_params": matches[1:],
					})

					for _, policyName := range route.Policies {
						policy, ok := module.Policies[policyName]
						if !ok {
							common.RespondInternalError(w)
							return
						}
						if !policy(*s.Cosys, r.WithContext(ctx)) {
							common.RespondError(w, "Forbidden", http.StatusForbidden)
							return
						}
					}

					uidRegex := regexp.MustCompile(`(.+)\.(.+)`)
					actionUid := route.Action
					uidMatches := uidRegex.FindStringSubmatch(actionUid)
					if len(uidMatches) != 3 {
						common.RespondInternalError(w)
						return
					}
					controllerName := uidMatches[1]
					actionName := uidMatches[2]
					controller, ok := module.Controllers[controllerName]
					if !ok {
						common.RespondInternalError(w)
						return
					}
					actionFunc, ok := controller[actionName]
					if !ok {
						common.RespondInternalError(w)
						return
					}
					action := actionFunc(*s.Cosys)

					for _, middlewareName := range route.Middlewares {
						middlewareFunc, ok := module.Middlewares[middlewareName]
						if !ok {
							common.RespondInternalError(w)
							return
						}
						middleware := middlewareFunc(*s.Cosys)

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
