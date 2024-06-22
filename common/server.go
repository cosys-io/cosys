package common

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var (
	svMutex    sync.RWMutex
	svRegister = make(map[string]func(*Cosys) Server)
)

func RegisterServer(svName string, svCtor func(*Cosys) Server) error {
	svMutex.Lock()
	defer svMutex.Unlock()

	if svCtor == nil {
		return fmt.Errorf("server is nil: %s", svName)
	}

	if _, dup := svRegister[svName]; dup {
		return fmt.Errorf("duplicate server:" + svName)
	}

	svRegister[svName] = svCtor
	return nil
}

type Server interface {
	Start() error
}

type ResponseCtxKey struct{}

type StateCtxKey struct{}

var (
	ResponseKey ResponseCtxKey
	StateKey    StateCtxKey
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

	states := ctx.Value(StateKey)
	if states == nil {
		return nil, fmt.Errorf("state map not found")
	}

	switch states.(type) {
	case map[string]any:

	default:
		return nil, fmt.Errorf("state has wrong type")
	}

	state, ok := states.(map[string]any)[stateName]
	if !ok {
		return nil, fmt.Errorf("state not found: %s", stateName)
	}

	return state, nil
}

func ReadParams(r *http.Request) (map[string]string, error) {
	params, err := ReadState(r, "queryParams")
	if err != nil {
		return nil, err
	}

	if params == nil {
		return nil, fmt.Errorf("query params not found")
	}

	switch params.(type) {
	case map[string]string:

	default:
		return nil, fmt.Errorf("query params has wrong type")
	}

	return params.(map[string]string), nil
}
