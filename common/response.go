package common

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

type Response struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func RespondOne(w http.ResponseWriter, data any, code int) {
	if w == nil {
		RespondInternalError(w)
		return
	}

	resp := Response{
		Data: data,
		Meta: Meta{},
	}

	header := w.Header()
	if header == nil {
		RespondInternalError(w)
		return
	}

	header.Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondInternalError(w)
	}
}

func RespondMany(w http.ResponseWriter, data any, page int, code int) {
	if w == nil {
		RespondInternalError(w)
		return
	}

	var size int
	if reflect.TypeOf(data).Kind() != reflect.Slice && reflect.TypeOf(data).Kind() != reflect.Array {
		size = reflect.ValueOf(data).Len()
	} else {
		size = 1
	}

	meta := Meta{}
	if page != 0 {
		meta.Pagination = &Pagination{
			Page:     page,
			PageSize: size,
		}
	}

	if data == nil {
		data = []any{}
	}

	resp := Response{
		Data: data,
		Meta: meta,
	}

	header := w.Header()
	if header == nil {
		RespondInternalError(w)
		return
	}

	header.Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondInternalError(w)
	}
}

func RespondError(w http.ResponseWriter, message string, code int) {
	if w == nil {
		RespondInternalError(w)
		return
	}

	resp := Response{
		Data: nil,
		Meta: Meta{
			Error: message,
		},
	}

	header := w.Header()
	if header == nil {
		RespondInternalError(w)
		return
	}

	header.Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondInternalError(w)
	}
}

func RespondInternalError(w http.ResponseWriter) {
	if w == nil {
		log.Printf("response writer is nil")
		return
	}

	resp := Response{
		Data: nil,
		Meta: Meta{
			Error: http.StatusText(http.StatusInternalServerError),
		},
	}

	header := w.Header()
	if header == nil {
		log.Printf("response header is nil")
		return
	}

	header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}
