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
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func RespondOne(w http.ResponseWriter, data any, code int) {
	resp := Response{
		Data: data,
		Meta: Meta{
			Pagination: Pagination{
				Page:     1,
				PageSize: 1,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Print(err)
	}
}

func RespondMany(w http.ResponseWriter, data any, page int, code int) {
	var size int
	if reflect.TypeOf(data).Kind() != reflect.Slice && reflect.TypeOf(data).Kind() != reflect.Array {
		size = reflect.ValueOf(data).Len()
	} else {
		size = 1
	}

	resp := Response{
		Data: data,
		Meta: Meta{
			Pagination: Pagination{
				Page:     page,
				PageSize: size,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Print(err)
	}
}

func RespondError(w http.ResponseWriter, message string, code int) {
	resp := Response{
		Data: message,
		Meta: Meta{
			Pagination: Pagination{
				Page:     1,
				PageSize: 1,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Print(err)
	}
}

func RespondInternalError(w http.ResponseWriter) {
	RespondError(w, "An Unexpected Error Has Occured.", http.StatusInternalServerError)
}
