package ws

import (
	"fmt"
	"log"

	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

func notifyStatusChange(statusType string, userID int) {
	connMu.Lock()
	defer func(){connMu.Unlock();println("f5321")}()
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

	var followerID int
	for rows.Next() {
		if err := rows.Scan(&followerID); err != nil {
			log.Println("Error scanning followerId:", err)
			continue
		}

		if followerConns, ok := Connections[followerID]; ok {
			msg := Message{
				Type:       "Status",
				StatusType: statusType,
				UserId:     userID,
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
	connMu.Lock()

	conns := Connections[userID]
	for i, c := range conns {
		if c == conn {
			Connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	connMu.Unlock()

	if len(Connections[userID]) == 0 {
		delete(Connections, userID)

		notifyStatusChange("offline", userID)

	}

	fmt.Println("loggs out successfully")
}
