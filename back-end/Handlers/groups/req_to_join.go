package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	"socialN/Handlers/ws"
	dataB "socialN/dataBase"
)

type ReqJoinGroup struct {
	GroupID int `json:"groupId"`
}

func Req_To_Join_Groups(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var req ReqJoinGroup
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if req.GroupID < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	var adminId int
	query := `SELECT adminId FROM Groups WHERE id = ?`
	err = dataB.SocialDB.QueryRow(query, req.GroupID).Scan(&adminId)
	if err != nil {
		fmt.Println(err)
		log.Println("Error to select admin in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	query2 := `
	   INSERT INTO Notifications (senderId, receiverId, type , groupId) VALUES (?,?,?,?)
	`

	res, err := dataB.SocialDB.Exec(query2, userID, adminId, "group_join_request", req.GroupID)
	if err != nil {
		fmt.Println(err)
		log.Println("Error to select admin in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}
	notifId, err := res.LastInsertId()
	if err != nil {
		fmt.Println("failed to extract notification I d")
		http.Error(w, "failed to extract notification I d", http.StatusInternalServerError)
		return
	}
	var firstName, lastName, avatar, groupName string
	err = dataB.SocialDB.QueryRow(`
	    SELECT u.firstName, u.lastName, u.avatar, g.title
	    FROM Users u
	    JOIN Groups g ON g.id = ?
	    WHERE u.id = ?
    `, req.GroupID, userID).Scan(&firstName, &lastName, &avatar, &groupName)
	if err != nil {
		log.Println("Error fetching sender and group info:", err)
		http.Error(w, "Could not get sender/group info", http.StatusInternalServerError)
		return
	}

	ws.ConnMu.Lock()
	conns := ws.Connections[adminId]
	ws.ConnMu.Unlock()

	for _, conn := range conns {
		if conn != nil {
			SendMessage(conn, InviteNotification{
				Id:               notifId,
				Type:             "Notification",
				NotificationType: "invite",
				GroupName:        groupName,
				FirstName:        firstName,
				LastName:         lastName,
				Avatar:           avatar,
			})
		}
	}
	w.WriteHeader(http.StatusAccepted)
}