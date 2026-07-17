package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"back-end/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	NickName  string `json:"nickName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
}
type checkRequest struct {
	Error error
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var data RegisterRequest
	w.Header().Set("content-type", "application/json")

	if r.Method != "POST" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "JSON Invalide")
		return
	}

	nickName := strings.TrimSpace(data.NickName)
	firstName := strings.TrimSpace(data.FirstName)
	lastName := strings.TrimSpace(data.LastName)
	email := strings.TrimSpace(data.Email)
	gender := strings.TrimSpace(data.Gender)
	password := strings.TrimSpace(data.Password)
	age := data.Age

	check, errMsg := ValidateRegistrationInput(nickName, firstName, lastName, email, gender, password, age)

	if !check {
		SendJSONResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if errHash != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user := models.User{
		NickName:  nickName,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
		Gender:    gender,
		Password:  string(hashedPassword),
	}

	if _, er := models.InsertUser(user); er != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	SendJSONResponse(w, http.StatusCreated, "User registered successfully")
}

func ValidateRegistrationInput(nickName, firstName, lastName, email, gender, password string, age int) (bool, string) {
	nickNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,10}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]{2,30}$`)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !nickNameRegex.MatchString(nickName) {
		return false, "Nickname must contain 3-10 letters, numbers or '_'"
	}

	exists, err := models.ExistsInColumn("nickName", nickName)
	if err != nil {
		return false, "Internal server error"
	}
	if exists {
		return false, "Nickname already exists"
	}

	if !nameRegex.MatchString(firstName) {
		return false, "Invalid first name"
	}

	if !nameRegex.MatchString(lastName) {
		return false, "Invalid last name"
	}

	if !emailRegex.MatchString(email) {
		return false, "Invalid email address"
	}
	exists, err = models.ExistsInColumn("email", email)
	if err != nil {
		return false, "Internal server error"
	}
	if exists {
		return false, "Email already exists"
	}

	if age < 13 || age > 120 {
		return false, "Invalid age"
	}

	if gender != "male" && gender != "female" {
		return false, "Gender must be 'male' or 'female'"
	}

	if len(password) < 8 {
		return false, "Password must contain at least 8 characters"
	}

	return true, ""
}