package server

import (
	"bytes"
	"context"
	"github.com/cosys-io/cosys/common"
	"net/http"
	"strings"
)

type Server struct {
	Port  string
	Cosys *common.Cosys
}

func (s Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api := s.Cosys.Api
		for _, route := range api.Routes {
			isMatch, params := matchPattern(r.URL.Path, route.Path)
			if isMatch {
				if r.Method != route.Method {
					continue
				}

				var log bytes.Buffer
				w = common.ResponseWriter{
					Writer: w,
					Log:    log,
				}

				queryParams := r.URL.Query()
				for paramName, param := range queryParams {
					if _, ok := params[paramName]; ok {
						common.RespondInternalError(w)
						return
					}
					if len(param) > 0 {
						params[paramName] = param[0]
					}
				}

				ctx := context.WithValue(r.Context(), common.ResponseKey, &log)
				ctx = context.WithValue(ctx, common.StateKey, map[string]any{
					"queryParams": params,
				})

				for _, policyName := range route.Policies {
					policy, ok := api.Policies[policyName]
					if !ok {
						common.RespondInternalError(w)
						return
					}
					if !policy(*s.Cosys, r.WithContext(ctx)) {
						common.RespondError(w, "Forbidden", http.StatusForbidden)
						return
					}
				}

				actionUid := route.Action
				actionUidSplit := strings.Split(actionUid, ".")
				if len(actionUidSplit) != 2 {
					common.RespondInternalError(w)
					return
				}
				controllerName := actionUidSplit[0]
				actionName := actionUidSplit[1]
				controller, ok := api.Controllers[controllerName]
				if !ok {
					common.RespondInternalError(w)
					return
				}
				actionCtor, ok := controller[actionName]
				if !ok {
					common.RespondInternalError(w)
					return
				}
				action := actionCtor(*s.Cosys)

				for _, middlewareName := range route.Middlewares {
					middlewareCtor, ok := api.Middlewares[middlewareName]
					if !ok {
						common.RespondInternalError(w)
						return
					}
					middleware := middlewareCtor(*s.Cosys)

					action = middleware(action)
				}

				action(w, r.WithContext(ctx))
				return
			}
		}
		http.NotFound(w, r)
	})

	if err := http.ListenAndServe(":"+s.Port, mux); err != nil {
		return err
	}

	return nil
}

func matchPattern(path, pattern string) (bool, map[string]string) {
	params := map[string]string{}

	questionIndex := strings.IndexByte(pattern, '?')
	if questionIndex != -1 {
		path = pattern[:questionIndex]
	}

	if len(path) != 0 && path[0] == '/' {
		path = path[1:]
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if len(pattern) != 0 && pattern[0] == '/' {
		pattern = pattern[1:]
	}
	if pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '{':
			paramNameEnd := strings.IndexByte(pattern, '}')
			paramName := pattern[1:paramNameEnd]
			pattern = pattern[paramNameEnd:]

			paramValueEnd := strings.IndexByte(path, '/')
			if paramValueEnd == -1 {
				paramValueEnd = len(path)
			}
			paramValue := path[:paramValueEnd]
			path = path[paramValueEnd:]

			params[paramName] = paramValue
		case path[0]:
			path = path[1:]
		default:
			return false, nil
		}
	}
	if path == "" && pattern == "" {
		return true, params
	} else {
		return false, nil
	}
}
