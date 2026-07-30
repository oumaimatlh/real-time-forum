package controllers

import (
	"net/http"
	"strconv"

	"back-end/middleware"
	"back-end/models"
)

func GetConversationMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	senderId := r.Context().Value(middleware.UserIdKey).(int)

	receiverIdStr := r.URL.Query().Get("receiverId")

	if receiverIdStr == "" {
		SendJSONResponse(w, http.StatusBadRequest, "receiverId is required", nil)
		return
	}

	receiverId, err := strconv.Atoi(receiverIdStr)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, "Invalid receiverId", nil)
		return
	}

	// pagination d messagees (historiques )
	limit := 10
	offset := 0

	if value := r.URL.Query().Get("limit"); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil {
			SendJSONResponse(
				w,
				http.StatusBadRequest,
				"Invalid limit",
				nil,
			)
			return
		}
	}

	if value := r.URL.Query().Get("offset"); value != "" {
		offset, err = strconv.Atoi(value)
		if err != nil {
			SendJSONResponse(
				w,
				http.StatusBadRequest,
				"Invalid offset",
				nil,
			)
			return
		}
	}

	messages, err := models.GetConversationMessages(
		senderId,
		receiverId,
		limit,
		offset,
	)
	if err != nil {
		SendJSONResponse(
			w,
			http.StatusInternalServerError,
			"Failed to get messages",
			nil,
		)
		return
	}

	SendJSONResponse(
		w,
		http.StatusOK,
		"Messages retrieved successfully",
		messages,
	)
}
