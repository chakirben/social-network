package ws

import (
	"encoding/json"
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
type goupmsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Sender  int    `json:"sender"`
	GroupID int    `json:"groupID"`
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

func OpenWsConn(resp http.ResponseWriter, req *http.Request) {
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
	defer deleteActiveuser(userID)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			deleteActiveuser(userID)
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
			err = MsgToDatabase(msg)
			if err == nil {
				Sendmessage(msg)
			} else {
				log.Println("Error sending message:", err)
			}

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
		case "groupmsg":
			var groupmsg goupmsg
			err := conn.ReadJSON(&groupmsg)
			if err != nil {
				log.Println("Error reading JSON:", err)
			} else {
				GroupMsgToDatabase(groupmsg)
				SendGroupMessage(groupmsg)
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

func removeConnection(username string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	conns := clients[username]
	for i, c := range conns {
		if c == conn {
			c.Close()
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	if len(conns) == 0 {
		delete(clients, username)
		msg := StatusChangeMessage{
			MessageType: "statusChange",
			UserName:    username,
			IsOnline:    false,
		}
		broadcastToAll(msg)
	} else {
		clients[username] = conns
	}
}

func deleteActiveuser(userID int) {
	fmt.Println("Deleting user:", userID)
	for i, conn := range Connections[userID] {
		if conn == nil {
			Connections[userID] = append(Connections[userID][:i], Connections[userID][i+1:]...)
			if len(Connections[userID]) == 0 {
				delete(Connections, userID)
			}
			break
		}
	}
	sendStutus(Status{
		Type:       "Status",
		StatusType: "offline",
		Receiver:   userID,
		Sender:     userID,
	})
}

func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	ids := []int{}
	for id := range Connections {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	query := "SELECT id, firstName, lastName, avatar FROM Users WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := dataB.SocialDB.Query(query, args...)
	if err != nil {
		http.Error(w, "Error", 500)
		return
	}
	defer rows.Close()

	type User struct {
		ID        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Avatar    string `json:"avatar"`
		Status    string `json:"status"`
	}

	var onlineUsers []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar)
		u.Status = "online"
		onlineUsers = append(onlineUsers, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(onlineUsers)
}
