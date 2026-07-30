package models

import (
	"time"

	"back-end/database"
)

type User struct {
	Id        int
	NickName  string
	FirstName string
	LastName  string
	Email     string
	Age       int
	Gender    string
	Password  string
	CreatedAt time.Time
}
type ConversationUser struct {
	Id            int
	NickName      string
	LastMessage   any
	LastMessageAt any
}

func InsertUser(user User) (int64, error) {
	query := "INSERT INTO users (nickName, firstName, lastName, email, Age, gender,  password) VALUES (?,?,?,?,?,?,?)"
	result, err := database.DB.Exec(query, user.NickName, user.FirstName, user.LastName, user.Email, user.Age, user.Gender, user.Password)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func GetUserByIdentifier(identifier string) (User, error) {
	user := User{}
	query := `
		SELECT id, nickName, firstName, lastName, email, Age, gender,  password, created_at
		FROM users
		WHERE nickName = ? OR email = ?
	`
	err := database.DB.QueryRow(query, identifier, identifier).Scan(&user.Id, &user.NickName, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.Gender, &user.Password, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {
	users := []User{}
	query := "SELECT * FROM users"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.NickName, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.Gender, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(userID int) (User, error) {
	user := User{}

	query := `
		SELECT id, nickName, firstName, lastName, email,
		       age, gender, password, created_at
		FROM users
		WHERE id = ?
	`

	err := database.DB.QueryRow(query, userID).Scan(
		&user.Id,
		&user.NickName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Age,
		&user.Gender,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ExistsInColumn(column string, value any) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users  WHERE " + column + " = ?"
	err := database.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func FilterConversationsUsers(currentId int) ([]ConversationUser, error) {
	users := []ConversationUser{}

	query := `
		SELECT
			u.id,
			u.nickName,

			(
				SELECT m.content
				FROM messages m
				WHERE
					(m.sender_id = ? AND m.receiver_id = u.id)
					OR
					(m.sender_id = u.id AND m.receiver_id = ?)
				ORDER BY m.created_at DESC, m.id DESC
				LIMIT 1
			) AS lastMessage,

			(
				SELECT m.created_at
				FROM messages m
				WHERE
					(m.sender_id = ? AND m.receiver_id = u.id)
					OR
					(m.sender_id = u.id AND m.receiver_id = ?)
				ORDER BY m.created_at DESC, m.id DESC
				LIMIT 1
			) AS lastMessageAt

		FROM users u

		WHERE u.id != ?

		ORDER BY
			CASE
				WHEN lastMessageAt IS NULL THEN 1
				ELSE 0
			END,

			lastMessageAt DESC,

			u.nickName ASC
	`

	rows, err := database.DB.Query(
		query,
		currentId,
		currentId,
		currentId,
		currentId,
		currentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := ConversationUser{}

		err := rows.Scan(
			&user.Id,
			&user.NickName,
			&user.LastMessage,
			&user.LastMessageAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
