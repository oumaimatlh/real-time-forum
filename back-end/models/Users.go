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

func ExistsInColumn(column, value string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users  WHERE " + column + " = ?"
	err := database.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
