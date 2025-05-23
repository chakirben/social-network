package ws

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"socialN/Handlers/auth"

	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Sender  int    `json:"sender"`
	Resever int    `json:"receiver"`
}
type Status struct {
	Type       string `json:"type"`
	StatusType string `json:"statusType"`
	Receiver   int    `json:"receiver"`
	Sender     int    `json:"sender"`
}
type Tyoping struct {
	Type       string `json:"type"`
	StatusType string `json:"statusType"`
	Sender     int    `json:"sender"`
	Receiver   int    `json:"receiver"`
}

var Connections = make(map[int][]*websocket.Conn)

func Entry(resp http.ResponseWriter, req *http.Request) {
	userID, err := auth.ValidateSession(req, dataB.SocialDB)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
	} else {
		Connections[userID] = append(Connections[userID], conn)
		fmt.Println("Connection upgraded successfully")
		sendStutus(Status{
			Type:       "Status",
			StatusType: "online",
			Receiver:   userID,
			Sender:     userID,
		})
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		msgtype := GetType(msg)
		switch msgtype {
		case "message":
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Error reading JSON:", err)
			}
			MsgToDatabase(msg)
			Sendmessage(msg)

		case "Status":
			var status Status
			err := conn.ReadJSON(&status)
			if err != nil {
				log.Println("Error reading JSON:", err)
			}
			sendStutus(status)
		case "typing":
			var typing Tyoping
			err := conn.ReadJSON(&typing)
			if err != nil {
				log.Println("Error reading JSON:", err)
				Typinsend(typing)
			}
		}
	}
}

func GetType(msg []byte) string {
	str := string(msg)
	str1 := strings.Split(str, ",")
	str1 = strings.Split(str1[0], ":")
	return strings.Trim(str1[1], "\"")
}
