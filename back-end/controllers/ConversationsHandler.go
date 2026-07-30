package controllers

import (
	"fmt"
	"net/http"

	"back-end/middleware"
	"back-end/models"
)

func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}
	userId := r.Context().Value(middleware.UserIdKey).(int)

	Users, err := models.FilterConversationsUsers(userId)
	if err != nil {
		fmt.Println(err.Error())
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	SendJSONResponse(w, http.StatusOK, "Conversations retrieved successfully.", Users)
}
