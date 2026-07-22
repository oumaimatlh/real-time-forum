package controllers

import (
	"encoding/json"
	"net/http"
)
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Message: message,
		Data:    data,
	})
}
