package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/Handlers/chat/models"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatHAndler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// websocketHandler(w, r)
	// Read messages and  insert it in the database
	msj := models.Message{}
	json.NewDecoder(r.Body).Decode(&msj)
	fmt.Println("message", msj)
	err := models.InsertMessage(models.Message{
		SenderID:   msj.SenderID,
		ReceiverID: msj.ReceiverID,
		Content:    msj.Content,
		CreatedAt:  msj.CreatedAt,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer upgrader.Close()
}
