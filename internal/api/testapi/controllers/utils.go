package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type Error struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

func WriteJSON(w http.ResponseWriter, data any, code int) {
	resp := Response{
		Code: code,
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		InternalErrorHandler(w)
		log.Print(err)
	}
}

func WriteError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		InternalErrorHandler(w)
		log.Print(err)
	}
}

func InternalErrorHandler(w http.ResponseWriter) {
	WriteError(w, "An Unexpected Error Has Occured.", http.StatusInternalServerError)
}
