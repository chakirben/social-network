package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Important: check origin manually
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// Replace with your frontend origin
		allowedOrigin := "http://localhost:3000"
		return origin == allowedOrigin
	},
}

type Message struct {
	Type string `json:"type"`

	// Common fields
	Content string `json:"content,omitempty"`
	Sender  int    `json:"sender,omitempty"`

	// For direct messages
	Receiver int `json:"receiver,omitempty"`

	// For group messages
	GroupID int `json:"groupID,omitempty"`

	// For status updates
	StatusType string `json:"statusType,omitempty"`
	UserId     int    `json:"userId,omitempty"`
}

var (
	Connections = make(map[int][]*websocket.Conn)
	connMu      sync.Mutex
)

func OpenWsConn(resp http.ResponseWriter, req *http.Request) {
	userID, err := auth.ValidateSession(req, dataB.SocialDB)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	connMu.Lock()

	Connections[userID] = append(Connections[userID], conn)

	connMu.Unlock()

	fmt.Println("Connection upgraded successfully")
	notifyStatusChange("online", userID)
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error reading JSON:", err)
			removeConn(userID, conn)
			break
		}
		fmt.Println("message : ", msg, Connections)
		switch msg.Type {
		case "message":
			msg.Sender = userID
			if err := RedirectMessage(msg); err != nil {
				log.Println(err)
				continue
			}
		case "groupmsg":
			if err := RedirectGroupMessage(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}
	println("conn closed ")
}
