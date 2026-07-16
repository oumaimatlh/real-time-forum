package models

import (
	"time"

	"back-end/database"
)

type Message struct {
	Id         int
	SenderId   int
	ReceiverId int
	Content    string
	CreatedAt  time.Time
}

func InsertMessage(message Message) error {
	query := "INSERT INTO messages (sender_id, receiver_id, content) VALUES (?,?,?)"
	_, err := database.DB.Exec(query, message.SenderId, message.ReceiverId, message.Content)
	if err != nil {
		return err
	}
	return nil
}

func GetConversationMessages(senderId, receiverId, limit, offset int) ([]Message, error) {
	messages := []Message{}
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE
			(sender_id = ? AND receiver_id = ?)
			OR
			(sender_id = ? AND receiver_id = ?)
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?;
	`
	rows, err := database.DB.Query(query, senderId, receiverId, receiverId, senderId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := Message{}
		err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}
