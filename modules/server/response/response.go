package response

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

// Response is the standard server response.
type Response struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

// Meta contains meta data about the response.
type Meta struct {
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination data about the response.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// RespondOne responds with no pagination data.
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

// RespondMany responds with pagination data.
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

	resp := Response{
		Data: data,
		Meta: Meta{
			Pagination: &Pagination{
				Page:     page,
				PageSize: size,
			},
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

// RespondError responds with an error message.
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

// RespondInternalError responds with an internal server error.
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
