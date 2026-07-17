package controllers

import (
	"encoding/json"
	"net/http"
)
func SendJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Message: message,
	})
}