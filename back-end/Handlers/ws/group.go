package ws

import (
	"fmt"
	"log"

	dataB "socialN/dataBase"
)

func RedirectGroupMessage(msg Message) error {
	// Check if sender is a member of the group
	println(msg.Sender , msg.GroupID)
	var isMember bool
	err := dataB.SocialDB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM GroupsMembers 
			WHERE memberId = ? AND groupId = ?
		)
	`, msg.Sender, msg.GroupID).Scan(&isMember)
	if err != nil {
		return fmt.Errorf("error checking group membership: %v", err)
	}
	if !isMember {
		return fmt.Errorf("sender is not a member of group %d", msg.GroupID)
	}

	// Save the group message
	_, err = dataB.SocialDB.Exec(`
		INSERT INTO GroupMessages (senderId, groupId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.GroupID, msg.Content)
	if err != nil {
		return fmt.Errorf("error inserting message into database: %v", err)
	}

	// Get all members of the group (excluding sender)
	rows, err := dataB.SocialDB.Query(`
		SELECT memberId FROM GroupsMembers 
		WHERE groupId = ?
	`, msg.GroupID, msg.Sender)
	if err != nil {
		return fmt.Errorf("error retrieving group members: %v", err)
	}
	defer rows.Close()

	// Send the message to all connected members
	for rows.Next() {
		var memberId int
		if err := rows.Scan(&memberId); err != nil {
			log.Println("Error scanning group member ID:", err)
			continue
		}
		println(Connections,memberId)
		for _, conn := range Connections[memberId] {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending group message to user %d: %v\n", memberId, err)
			}
		}
	}

	return nil
}