package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"back-end/models"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data LoginRequest

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "JSON Invalid")
		return
	}

	data.Identifier = strings.TrimSpace(data.Identifier)
	if data.Identifier == "" {
		SendJSONResponse(w, http.StatusBadRequest, "Identifier is required")
		return
	}
	if data.Password == "" {
		SendJSONResponse(w, http.StatusBadRequest, "Password is required")
		return
	}
	user, err := models.GetUserByIdentifier(data.Identifier)
	if err != nil {
		SendJSONResponse(w, http.StatusUnauthorized, "No Account with this identifier")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		SendJSONResponse(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	err = models.DeleteSessionsByUserID(user.Id)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	token, err := models.InsertSession(user.Id)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	SendJSONResponse(w, http.StatusOK, "Login successful")
}
