package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type Notification struct {
	ID               int       `json:"id"`
	SenderFirstName  string    `json:"senderFirstName"`
	SenderLastName   string    `json:"senderLastName"`
	SenderAvatar     string    `json:"senderAvatar"`
	ReceiverID       *int      `json:"receiverId,omitempty"`
	Type             string    `json:"type"`
	Status           string    `json:"status"`
	NotificationDate time.Time `json:"notificationDate"`
	GroupID          *int      `json:"groupId,omitempty"`
	EventID          *int      `json:"eventId,omitempty"`
	GroupTitle       *string   `json:"groupTitle,omitempty"`
	EventTitle       *string   `json:"eventTitle,omitempty"`
}

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	// Parse userId from query
	userId, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	query := `
	SELECT 
		n.id,
		u.firstName,
		u.lastName,
		u.avatar,
		n.receiverId,
		n.type,
		n.status,
		n.notificationDate,
		n.groupId,
		g.title AS groupTitle,
		n.eventId,
		e.title AS eventTitle
	FROM 
		Notifications n
	JOIN 
		Users u ON u.id = n.senderId
	LEFT JOIN 
		Groups g ON g.id = n.groupId
	LEFT JOIN
		Events e ON e.id = n.eventId
	WHERE 
		n.receiverId = ?
	ORDER BY 
		n.notificationDate DESC;
`

	rows, err := dataB.SocialDB.Query(query, userId)
	if err != nil {
		http.Error(w, "Failed to query notifications", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID,
			&n.SenderFirstName,
			&n.SenderLastName,
			&n.SenderAvatar,
			&n.ReceiverID,
			&n.Type,
			&n.Status,
			&n.NotificationDate,
			&n.GroupID,
			&n.GroupTitle,
			&n.EventID,
			&n.EventTitle,
		)
		if err != nil {
			fmt.Println("Error scanning notification row:", err)
			http.Error(w, "Error scanning notification row", http.StatusInternalServerError)
			return
		}
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error finalizing row iteration", http.StatusInternalServerError)
		return
	}

	// Send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
