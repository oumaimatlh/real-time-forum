package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"back-end/models"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	posts, err := models.GetAllPosts()
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)


	fmt.Println(posts)
}
