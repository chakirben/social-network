package ws

import (
	"fmt"

	dataB "socialN/dataBase"
)

func Sendmessage(msg Message) {
	for _, conn := range Connections[msg.Resever] {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

func MsgToDatabase(msg Message) error{
	var follows bool
	err := dataB.SocialDB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Follows 
			WHERE (followerId = ? AND followingId = ?)
			   OR (followerId = ? AND followingId = ?)
		)`, msg.Sender, msg.Resever, msg.Resever, msg.Sender).Scan(&follows)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %v", err)
	}
	if !follows {
		return fmt.Errorf("sender and receiver do not follow each other")
	}
	_, err = dataB.SocialDB.Exec(`
		INSERT INTO Comments (senderId, receiverId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.Resever, msg.Content)
	if err != nil {
		return fmt.Errorf("error inserting message into database: %v", err)
	}
	return nil
}
