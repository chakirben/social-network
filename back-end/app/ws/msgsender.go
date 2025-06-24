package ws

import (
	"fmt"

	dataB "socialN/dataBase"
)

func RedirectMessage(msg Message) error {
	allowed := checkAuthorisation(msg)
	if !allowed {
		return fmt.Errorf("Not allowed to msg this user")
	}
	_, err := dataB.SocialDB.Exec(`
		INSERT INTO Messages (senderId, receiverId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.Receiver, msg.Content)
	if err != nil {
		return fmt.Errorf("error inserting message into database: %v", err)
	}

	err = dataB.SocialDB.QueryRow(`
		SELECT firstName, lastName, avatar FROM Users WHERE id = ?
	`, msg.Sender).Scan(&msg.FirstName, &msg.LastName, &msg.Avatar)
	if err != nil {
		return fmt.Errorf("error fetching user details: %v", err)
	}
	for _, conn := range Connections[msg.Receiver] {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
	for _, conn := range Connections[msg.Sender] {
		println("sent it to me")
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
	return nil
}

func checkAuthorisation(msg Message) bool {
	fmt.Println("Checking authorization: sender =", msg.Sender, "receiver =", msg.Receiver)

	var isFollowing bool
	err := dataB.SocialDB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM Followers 
			WHERE (followerId = ? AND followedId = ?)
			   OR (followerId = ? AND followedId = ?)
		)
	`, msg.Sender, msg.Receiver, msg.Receiver, msg.Sender).Scan(&isFollowing)
	if err != nil {
		fmt.Println("Error checking follower relationship:", err)
		return false
	}

	return isFollowing
}
