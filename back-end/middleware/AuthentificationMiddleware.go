package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"back-end/models"
)

type contextKey string

const UserIdKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			writeUnauthorizedResponse(w)
			return
		}
		session, err := models.GetSessionByToken(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
				Path:     "/",
			})
			writeUnauthorizedResponse(w)
			return
		}
		cts := context.WithValue(r.Context(), UserIdKey, session.UserId)
		next(w, r.WithContext(cts))
	}
}

func writeUnauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Unauthorized",
		"data":    nil,
	})
}
