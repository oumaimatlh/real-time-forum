package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"back-end/models"

	"github.com/gorilla/websocket"
)

var clients map[int]*websocket.Conn

type MessagePayload struct { // hade la structure li ghade tseftahali mn front end
	ReceiverId int    `json:"receiverId"`
	Message    string `json:"message"`
}

// Transformer une connexion HTTP en connexion WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value("userId").(int)

	con, err := upgrader.Upgrade(w, r, nil) // obj connection webSoc grace a lui on peut lire les messages && envoyer ..
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, "Failed to upgrade to WebSocket", nil)
		return
	}

	clients[senderId] = con

	defer con.Close()

	for {
		messageType, data, err := con.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.TextMessage {
			con.WriteJSON(map[string]string{"errMsg": "Unsupported message type"})
			continue
		}

		payload := MessagePayload{}

		err = json.Unmarshal(data, &payload)
		if err != nil {
			con.WriteJSON(map[string]string{"errMsg": err.Error()})
			continue
		}

		existReceiver, err := models.GetUserByID(payload.ReceiverId)
		if err != nil {
			con.WriteJSON(map[string]string{"errMsg": "Receiver not found"})
			continue
		}
		if payload.ReceiverId == senderId {
			con.WriteJSON(map[string]string{"errMsg": "You cannot send a message to yourself"})
			continue
		}
		if payload.Message == ""  || len(payload.Message) > 500 {
			con.WriteJSON(map[string]string{"errMsg": "Invalid message"})
			continue
		}


		message := models.Message{
			SenderId:   senderId,
			ReceiverId: payload.ReceiverId,
			Content:    payload.Message,
		}
		err = models.InsertMessage(message)
		if err != nil {
			con.WriteJSON(map[string]string{"errMsg": "Failed to insert message"})
			continue
		}

		if receiverCon, ok := clients[payload.ReceiverId]; ok {
			err = receiverCon.WriteJSON(map[string]interface{}{
				"senderId": senderId,
				"message":  payload.Message,
			})
			if err != nil {
				fmt.Println("Failed to send message to receiver:", err)
			}
		}

	}
}