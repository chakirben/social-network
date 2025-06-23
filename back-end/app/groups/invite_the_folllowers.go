package groups

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"socialN/app/auth"
	"socialN/app/ws"
	dataB "socialN/dataBase"
)

type TheUserInvited struct {
	UserId  int `json:"userId"`
	GroupId int `json:"groupId"`
}

type InviteNotification struct {
	Id               int64  `json:"id"`
	Type             string `json:"type"`
	NotificationType string `json:"notificationType"`
	GroupName        string `json:"title"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Avatar           string `json:"avatar"`
}

func InviteTheFollowers(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session ", http.StatusUnauthorized)
		return
	}

	var userIDToInvite TheUserInvited
	if err := json.NewDecoder(r.Body).Decode(&userIDToInvite); err != nil {
		http.Error(w, "Invalid JSON ", http.StatusBadRequest)
		return
	}

	// Check if the invited user exists
	var exists bool
	err = dataB.SocialDB.QueryRow(`SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)`, userIDToInvite.UserId).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Insert notification
	query := `
		INSERT INTO Notifications (senderId, receiverId, type, groupId, status)
		VALUES (?, ?, ?, ?, 'pending')
	`
	result, err := dataB.SocialDB.Exec(query, userID, userIDToInvite.UserId, "group_invite", userIDToInvite.GroupId)
	if err != nil {
		http.Error(w, "Could not send invitation", http.StatusInternalServerError)
		return
	}
	notifId, _ := result.LastInsertId()

	// Fetch sender (inviter) info
	var firstName, lastName, avatar string
	err = dataB.SocialDB.QueryRow(`
		SELECT firstName, lastName, avatar FROM Users WHERE id = ?
	`, userID).Scan(&firstName, &lastName, &avatar)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error fetching sender info", http.StatusInternalServerError)
		return
	}

	// Fetch group name
	var groupName string
	err = dataB.SocialDB.QueryRow(`SELECT title FROM Groups WHERE id = ?`, userIDToInvite.GroupId).Scan(&groupName)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error fetching group name", http.StatusInternalServerError)
		return
	}

	// Send WebSocket notification
	ws.ConnMu.Lock()
	conns := ws.Connections[userIDToInvite.UserId]
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
