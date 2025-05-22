package ws

import (
	"fmt"

	dataB "socialN/dataBase"
)

func sendStutus(msg Status) {
	for _, users := range connections {
		for _, conn1 := range users {
			if checkAlreadyFollow(msg.Receiver, msg.Sender) { // this condition's not working because i dint know how to check if the user is follow by this user or not
				err := conn1.WriteJSON(msg)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
			}
		}
	}
}

func checkAlreadyFollow(followerID, followedID int) bool {
	rows, err := dataB.SocialDB.Query(
		`SELECT 1 FROM Followers WHERE followerId = ? AND followedId = ?`,
		followerID, followedID,
	)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}
