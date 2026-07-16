package models

import (
	"back-end/database"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	IdSession int
	UserId    int
	Token     string
	ExpiresAt time.Time
}

func InsertSession(idUser int) (string, error) {
	query := "INSERT INTO session (user_id, token, expires_at) VALUES (?, ?, ?)"
	token := uuid.Must(uuid.NewV4()).String()
	expiresAt := time.Now().Add(time.Hour)
	_, err := database.DB.Exec(query, idUser, token, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetSessionByToken(token string) (Session, error) {
	session := Session{}
	query := "SELECT id, user_id, token, expires_at FROM session WHERE token = ?"
	err := database.DB.QueryRow(query, token).Scan(&session.IdSession, &session.UserId, &session.Token, &session.ExpiresAt)
	if err != nil {
		return Session{}, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		DeleteSessionByToken(token)
		return Session{}, errors.New("session expirée")
	}
	return session, nil
}

func DeleteSessionByToken(token string) error {
	query := "DELETE  FROM session WHERE token = ?"
	_, err := database.DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSessionsByUserID(userID int) error {
	_, err := database.DB.Exec("DELETE FROM session WHERE user_id = ?", userID)
	return err
}
