package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"back-end/middleware"
	"back-end/models"
)

type CommentRequest struct {
	UserId  int    `json:"userId"`
	Content string `json:"content"`
}

func CreateCommentPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}
	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid post id", nil)
		return
	}
	var data CommentRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	if err := ValidateCommentInput(data); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	_, err = models.GetPostByID(postID)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "Post not found", nil)
		return
	}
	_, err = models.ExistsInColumn("id", data.UserId)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}
	comment := models.Comments{
		UserId:  data.UserId,
		PostId:  postID,
		Content: data.Content,
	}
	err = models.InsertComment(comment)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	SendJSONResponse(w, http.StatusCreated, "Comment created successfully", nil)
}

func ValidateCommentInput(data CommentRequest) error {
	data.Content = strings.TrimSpace(data.Content)

	if data.Content == "" {
		return errors.New("content is required")
	}

	if len(data.Content) > 1000 {
		return errors.New("content must not exceed 1000 characters")
	}

	return nil
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	reactToComment(w, r, "like")
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	reactToComment(w, r, "dislike")
}

func reactToComment(w http.ResponseWriter, r *http.Request, reactionType string) {
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	commentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid comment id", nil)
		return
	}

	comment, err := models.GetCommentByID(commentID)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "Comment not found", nil)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)

	err = models.ToggleCommentReaction(
		userID,
		comment.PostId,
		commentID,
		reactionType,
	)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	likes, _ := models.CountLikesByComments(commentID, comment.PostId, "like")
	dislikes, _ := models.CountLikesByComments(commentID, comment.PostId, "dislike")

	SendJSONResponse(w, http.StatusOK, "Reaction updated", map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	})
}
