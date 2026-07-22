package controllers

import (
	"back-end/models"
	"net/http"
)

func GetCategoryHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		SendJSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}
	category, err := models.GetAllCategories()
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	SendJSONResponse(w, http.StatusOK, "Category retrieved successfully", category)
}