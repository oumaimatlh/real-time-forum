package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"back-end/middleware"
	"back-end/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

var (
	clients     = map[int][]*Client{}
	onlineUsers = map[int]bool{}
	mu          sync.RWMutex // protecte jusre ces 2 maps
)

type MessagePayload struct {
	ReceiverId int    `json:"receiverId"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type OnlineUsers struct {
	Type  string `json:"type"`
	Users []int  `json:"users"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	senderId := r.Context().Value(middleware.UserIdKey).(int)
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn: con,
	}

	mu.Lock()
	clients[senderId] = append(clients[senderId], client)
	onlineUsers[senderId] = true
	mu.Unlock()

	BroadCastOnlineUsers()

	defer func() {
		mu.Lock()
		userClients := clients[senderId]

		for i, connectedClient := range userClients {
			if connectedClient == client {
				clients[senderId] = append(
					userClients[:i],
					userClients[i+1:]...,
				)
				break
			}
		}
		if len(clients[senderId]) == 0 {
			delete(clients, senderId)
			delete(onlineUsers, senderId)
		}
		mu.Unlock()

		BroadCastOnlineUsers()

		con.Close()
	}()

	for {
		messageType, data, err := con.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.TextMessage {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Unsupported message type",
			})

			continue
		}

		payload := MessagePayload{}

		err = json.Unmarshal(data, &payload)
		if err != nil {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Invalid message format",
			})
			continue
		}

		if payload.ReceiverId <= 0 {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Invalid receiver id",
			})
			continue
		}

		_, err = models.GetUserByID(payload.ReceiverId)
		if err != nil {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Receiver not found",
			})

			continue
		}

		if payload.ReceiverId == senderId {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "You cannot send message to yourself",
			})
			continue
		}

		if strings.TrimSpace(payload.Message) == "" || len(payload.Message) > 500 {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Invalid message",
			})

			continue
		}

		message := models.Message{
			SenderId:   senderId,
			ReceiverId: payload.ReceiverId,
			Content:    payload.Message,
		}

		err = models.InsertMessage(message)
		if err != nil {
			client.Send(ErrorResponse{
				Type:    "error",
				Message: "Failed to insert message",
			})

			continue
		}

		mu.RLock()
		receiverClients := clients[payload.ReceiverId]
		mu.RUnlock()

		if len(receiverClients) > 0 {
			for _, receiver := range receiverClients {

				err = receiver.Send(map[string]interface{}{
					"type":     "private_message",
					"senderId": senderId,
					"message":  message,
				})
				if err != nil {
					fmt.Println("Failed sending message:", err)
				}
			}
		} else {
			client.Send(map[string]interface{}{
				"type":    "info",
				"message": "User is offline. Message saved.",
			})
		}
	}
}

func (c *Client) Send(data interface{}) error {
	c.Mu.Lock()

	defer c.Mu.Unlock()

	return c.Conn.WriteJSON(data)
}

func BroadCastOnlineUsers() {
	mu.RLock()

	users := make([]int, 0, len(onlineUsers))

	clientsCopy := make([]*Client, 0)

	for id := range onlineUsers {
		users = append(users, id)
	}

	for _, userClients := range clients {
		for _, client := range userClients {
			clientsCopy = append(clientsCopy, client)
		}
	}

	mu.RUnlock()

	payload := OnlineUsers{
		Type:  "onlineUsers",
		Users: users,
	}

	for _, client := range clientsCopy {

		err := client.Send(payload)
		if err != nil {
			fmt.Println("Broadcast error:", err)
		}
	}
}
