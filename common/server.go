package common

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var (
	svMutex sync.RWMutex
	svMap   = make(map[string]func(*Cosys) Server)
)

func RegisterServer(name string, server func(*Cosys) Server) error {
	svMutex.Lock()
	defer svMutex.Unlock()

	if server == nil {
		return fmt.Errorf("server is nil")
	}

	if _, dup := svMap[name]; dup {
		return fmt.Errorf("duplicate server:" + name)
	}

	svMap[name] = server
	return nil
}

type Server interface {
	Start() error
}

type ResponseContextKey struct{}

type StateContextKey struct{}

var (
	ResponseKey ResponseContextKey
	StateKey    StateContextKey
)

type ResponseWriter struct {
	Writer http.ResponseWriter
	Log    bytes.Buffer
}

func (r ResponseWriter) Header() http.Header {
	return r.Writer.Header()
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	r.Log.Write(b)
	return r.Writer.Write(b)
}

func (r ResponseWriter) WriteHeader(statusCode int) {
	r.Writer.WriteHeader(statusCode)
}

func ReadState(r *http.Request, stateName string) (any, error) {
	ctx := r.Context()
	if ctx == nil {
		return nil, fmt.Errorf("context not found")
	}

	stateMap := ctx.Value(StateKey)
	if stateMap == nil {
		return nil, fmt.Errorf("state not found")
	}

	switch stateMap.(type) {
	case map[string]any:

	default:
		return nil, fmt.Errorf("state has wrong type")
	}

	state, ok := stateMap.(map[string]any)[stateName]
	if !ok {
		return nil, fmt.Errorf("state not found: " + stateName)
	}

	return state, nil
}

func ReadParams(r *http.Request) ([]string, error) {
	params, err := ReadState(r, "query_params")
	if err != nil {
		return nil, err
	}

	if params == nil {
		return nil, fmt.Errorf("query params not found")
	}

	switch params.(type) {
	case []string:

	default:
		return nil, fmt.Errorf("query params has wrong type")
	}

	return params.([]string), nil
}
