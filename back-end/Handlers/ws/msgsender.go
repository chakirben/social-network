package ws

import (
	"fmt"

	dataB "socialN/dataBase"
)

func Sendmessage(msg Message) {
	for _, conn := range connections[msg.Resever] {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

func MsgToDatabase(msg Message) {
	_, err := dataB.SocialDB.Exec(`
		INSERT INTO Comments (senderId, receiverId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.Resever, msg.Content)
	if err != nil {
		fmt.Println("Error inserting message into database:", err)
	}
}
