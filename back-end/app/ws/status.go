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

	// Find followers
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

	// Find following
	queryFollowing := `
		SELECT followedId
		FROM Followers
		WHERE followerId = ?
	`
	rowsFollowing, err := dataB.SocialDB.Query(queryFollowing, userID)
	if err != nil {
		log.Println("Failed to query following:", err)
		return
	}
	defer rowsFollowing.Close()

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

	// Send to followers
	var followerID int
	for rows.Next() {
		if err := rows.Scan(&followerID); err != nil {
			log.Println("Error scanning followerId:", err)
			continue
		}
		if followerConns, ok := Connections[followerID]; ok {
			for _, conn := range followerConns {
				if err := conn.WriteJSON(msg); err != nil {
					log.Println("Error sending status to follower:", err)
				} else {
					fmt.Println("Status update sent to follower:", followerID)
				}
			}
		}
	}

	// Send to following
	var followingID int
	for rowsFollowing.Next() {
		if err := rowsFollowing.Scan(&followingID); err != nil {
			log.Println("Error scanning followingId:", err)
			continue
		}
		if followingConns, ok := Connections[followingID]; ok {
			for _, conn := range followingConns {
				if err := conn.WriteJSON(msg); err != nil {
					log.Println("Error sending status to following user:", err)
				} else {
					fmt.Println("Status update sent to following user:", followingID)
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
