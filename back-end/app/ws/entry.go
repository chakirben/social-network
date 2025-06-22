package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"socialN/app/auth"
	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// allowedOrigin := "http://127.0.0.1:3000"
		//return origin == allowedOrigin
		return true 
	},
}

type Message struct {
	Type       string `json:"type"`
	Content    string `json:"content,omitempty"`
	Sender     int    `json:"sender,omitempty"`
	Receiver   int    `json:"receiverId,omitempty"`
	GroupID    int    `json:"groupID,omitempty"`
	StatusType string `json:"statusType,omitempty"`
	UserId     int    `json:"userId,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
}

var (
	Connections = make(map[int][]*websocket.Conn)
	ConnMu      sync.Mutex
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

	ConnMu.Lock()

	Connections[userID] = append(Connections[userID], conn)

	ConnMu.Unlock()

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
			fmt.Println("message type is message")
			msg.Sender = userID
			if err := RedirectMessage(msg); err != nil {
				log.Println(err)
				continue
			}
		case "groupmsg":
			msg.Sender = userID
			if err := RedirectGroupMessage(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}
	println("conn closed ")
}
