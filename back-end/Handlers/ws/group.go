package ws

import (
	"fmt"
	dataB "socialN/dataBase"
)

func GroupMsgToDatabase(msg goupmsg) {
	_, err := dataB.SocialDB.Exec(`
		INSERT INTO GroupMessages (senderId, groupId, content)
		VALUES (?, ?, ?)`, msg.Sender, msg.GroupID, msg.Content)
	if err != nil {
		fmt.Println("Error inserting message into database:", err)
	}
}
func SendGroupMessage(msg goupmsg) {
	rows, err := dataB.SocialDB.Query(`
		SELECT memberId FROM GroupMembers WHERE groupId = ?`, msg.GroupID)
	if err != nil {
		fmt.Println("Error querying group members:", err)
		return
	}
	defer rows.Close()

	var memberIDs []int
	for rows.Next() {
		var memberID int
		if err := rows.Scan(&memberID); err != nil {
			fmt.Println("Error scanning memberId:", err)
			continue
		}
		memberIDs = append(memberIDs, memberID)
	}
	for _, memberID := range memberIDs {
		if memberID != msg.Sender {
			for _, conn := range Connections[memberID] {
				if err := conn.WriteJSON(msg); err != nil {
					fmt.Println("Error sending message to member:", err)
				}
			}
		}
	}

}
