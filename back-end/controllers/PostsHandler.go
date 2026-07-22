package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"back-end/models"
)

type PostRequest struct {
	UserId   int    `json:UserId`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category []int  `json:category`
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	posts, err := models.GetAllPosts()
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	SendJSONResponse(w, http.StatusOK, "Posts retrieved successfully.", posts)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var data PostRequest
	if r.Method != "POST" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "JSON Invalide", nil)
		return
	}
	if err := ValidatePostInput(data); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	_, err = models.ExistsInColumn("id", data.UserId)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	for _, categoryID := range data.Category {
		_, err := models.GetCategoryByID(categoryID)
		if err != nil {
			SendJSONResponse(w, http.StatusBadRequest, "Invalid category", nil)
			return
		}
	}
	post := models.Post{
		Title:   data.Title,
		Content: data.Content,
		UserId:  data.UserId,
	}

	_, err = models.InsertPost(post)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	SendJSONResponse(w, http.StatusCreated, "Post created successfully", nil)
}

func GetPostsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	categoryID := r.URL.Query().Get("category")
	if categoryID == "" {
		SendJSONResponse(w, http.StatusBadRequest, "Category is required", nil)
		return
	}

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid category", nil)
		return
	}

	posts, err := models.GetPostsByCategory(id)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	SendJSONResponse(w, http.StatusOK, "Posts retrieved successfully", posts)
}

func ValidatePostInput(data PostRequest) error {
	data.Title = strings.TrimSpace(data.Title)
	data.Content = strings.TrimSpace(data.Content)

	if data.Title == "" {
		return errors.New("title is required")
	}

	if len(data.Title) > 100 {
		return errors.New("title must not exceed 100 characters")
	}

	if data.Content == "" {
		return errors.New("content is required")
	}
	if len(data.Content) > 3000 {
		return errors.New("title must not exceed 100 characters")
	}

	if len(data.Category) == 0 {
		return errors.New("At least one category is required")
	}

	return nil
}
