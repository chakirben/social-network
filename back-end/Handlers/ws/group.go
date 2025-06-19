package ws

import (
	"database/sql"
	"fmt"
	"log"

	dataB "socialN/dataBase"
)


// RedirectGroupMessage saves a group message and sends it to all group members
func RedirectGroupMessage(msg Message) error {
	// 1. Check if sender is a member of the group
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

	// 2. Fetch user info
	var firstName, lastName, avatar sql.NullString
	err = dataB.SocialDB.QueryRow(`
		SELECT firstName, lastName, avatar FROM Users
		WHERE id = ?
	`, msg.Sender).Scan(&firstName, &lastName, &avatar)
	if err != nil {
		return fmt.Errorf("error fetching user info: %v", err)
	}
	msg.FirstName = firstName.String
	msg.LastName = lastName.String
	msg.Avatar = avatar.String

	// 3. Insert the group message into the database
	_, err = dataB.SocialDB.Exec(`
		INSERT INTO GroupMessages (senderId, groupId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.GroupID, msg.Content)
	if err != nil {
		return fmt.Errorf("error inserting message: %v", err)
	}

	// 4. Get all group members (excluding the sender)
	rows, err := dataB.SocialDB.Query(`
		SELECT memberId FROM GroupsMembers 
		WHERE groupId = ? AND memberId != ?
	`, msg.GroupID, msg.Sender)
	if err != nil {
		return fmt.Errorf("error retrieving group members: %v", err)
	}
	defer rows.Close()

	// 5. Send enriched message to all group members
	for rows.Next() {
		var memberId int
		if err := rows.Scan(&memberId); err != nil {
			log.Println("error scanning memberId:", err)
			continue
		}
		for _, conn := range Connections[memberId] {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("error sending message to user %d: %v\n", memberId, err)
			}
		}
	}

	return nil
}
