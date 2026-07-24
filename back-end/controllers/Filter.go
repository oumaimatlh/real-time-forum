package controllers

import (
	"net/http"
	"strings"

	"back-end/middleware"
	"back-end/models"
)

func FilterPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	query := r.URL.Query()

	// ?category=1,2,3
	categoryParam := query.Get("category")

	// ?mine=true
	filterMine := query.Get("mine") == "true"

	// ?liked=true
	filterLiked := query.Get("liked") == "true"

	// Aucun filtre
	if categoryParam == "" && !filterMine && !filterLiked {
		posts, err := models.GetAllPosts()
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		SendJSONResponse(w, http.StatusOK, "Posts found", posts)
		return
	}

	userID := 0

	if filterMine || filterLiked {
		id, ok := r.Context().Value(middleware.UserIdKey).(int)
		if !ok {
			SendJSONResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		userID = id
	}

	var categories []string

	if categoryParam != "" {
		categories = strings.Split(categoryParam, ",")
	}

	posts, err := models.GetFilteredPosts(userID, categories, filterLiked, filterMine)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	SendJSONResponse(w, http.StatusOK, "Posts found", posts)
}
