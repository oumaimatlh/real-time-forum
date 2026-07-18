package controllers

import (
	"net/http"
	"time"

	"back-end/models"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		SendJSONResponse(w, http.StatusUnauthorized, " Unauthorized")
		return
	}

	err = models.DeleteSessionByToken(cookie.Value)
	if err != nil {
		SendJSONResponse(w, http.StatusUnauthorized, " Invalid session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	SendJSONResponse(w, http.StatusOK, "Logout successful")
}
