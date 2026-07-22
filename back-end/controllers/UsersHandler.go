package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"back-end/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}
	users, err := models.GetAllUsers()
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	SendJSONResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/users/")
	id = strings.Trim(id, "/")
	if id == "" {
		SendJSONResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid User ID", nil)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	SendJSONResponse(w, http.StatusOK, "User retrieved successfully", user)
}
