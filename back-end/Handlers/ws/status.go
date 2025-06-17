package ws

import (
	"fmt"
	"log"

	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

// UserStatusPayload holds public user info for status updates
type UserStatusPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
}

// StatusMessage is sent to followers when a user's status changes
type StatusMessage struct {
	Type       string             `json:"type"`       // "Status"
	StatusType string             `json:"statusType"` // "online" or "offline"
	UserId     int                `json:"userId"`
	User       *UserStatusPayload `json:"user"`
}

func notifyStatusChange(statusType string, userID int) {
	ConnMu.Lock()
	defer func() { ConnMu.Unlock(); println("f5321") }()

	// Find followers of the user
	query := `
		SELECT followerId 
		FROM Followers 
		WHERE followedId = ?
	`
	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		log.Println("Failed to query followers:", err)
		return
	}
	defer rows.Close()

	// Fetch profile of the user whose status changed
	var firstName, lastName, avatar string
	userQuery := `
		SELECT firstName, lastName, avatar
		FROM Users
		WHERE id = ?
	`
	err = dataB.SocialDB.QueryRow(userQuery, userID).Scan(&firstName, &lastName, &avatar)
	if err != nil {
		log.Println("Failed to get user profile info:", err)
		return
	}

	// Send update to each follower
	var followerID int
	for rows.Next() {
		if err := rows.Scan(&followerID); err != nil {
			log.Println("Error scanning followerId:", err)
			continue
		}

		if followerConns, ok := Connections[followerID]; ok {
			msg := StatusMessage{
				Type:       "Status",
				StatusType: statusType,
				UserId:     userID,
				User: &UserStatusPayload{
					FirstName: firstName,
					LastName:  lastName,
					Avatar:    avatar,
				},
			}

			for _, conn := range followerConns {
				if err := conn.WriteJSON(msg); err != nil {
					log.Println("Error sending status to follower:", err)
				}
			}
		}
	}
}

func removeConn(userID int, conn *websocket.Conn) {
	ConnMu.Lock()

	conns := Connections[userID]
	for i, c := range conns {
		if c == conn {
			Connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	ConnMu.Unlock()

	if len(Connections[userID]) == 0 {
		delete(Connections, userID)

		notifyStatusChange("offline", userID)

	}

	fmt.Println("loggs out successfully")
}
