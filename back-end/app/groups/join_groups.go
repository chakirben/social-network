package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

type JoinGroupRequest struct {
	GroupID int `json:"groupId"`
	UserId  int `json:"userId"`
}

// Insert the GroupsMembers in the database....
func JoinGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var req JoinGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if req.GroupID < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO GroupsMembers (memberId, groupId) VALUES (?, ?)
	`
	_, err = dataB.SocialDB.Exec(query, req.UserId, req.GroupID)
	if err != nil {
		log.Println("Error to insert members in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func SendMessage(conn *websocket.Conn, msg any) {
	fmt.Println("Sending message to client:", msg)
	message, err := json.Marshal(msg)
	if err != nil {
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}