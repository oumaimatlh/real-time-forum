package controllers

import (
	"net/http"
	"strconv"

	"back-end/middleware"
	"back-end/models"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	reactToPost(w, r, "like")
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	reactToPost(w, r, "dislike")
}

func reactToPost(w http.ResponseWriter, r *http.Request, reactionType string) {
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid post id", nil)
		return
	}

	_, err = models.GetPostByID(postID)
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, "Post not found", nil)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(int)

	err = models.TogglePostReaction(userID, postID, reactionType)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	likes, _ := models.CountLikeDislikeByPost(postID, "like")
	dislikes, _ := models.CountLikeDislikeByPost(postID, "dislike")

	SendJSONResponse(w, http.StatusOK, "Reaction updated", map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	})
}
