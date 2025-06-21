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
		n.eventId
	FROM 
		Notifications n
	JOIN 
		Users u ON u.id = n.senderId
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
			&n.EventID,
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

// type loggedUser struct {
// 	UserSession string `json:"user_session"`
// }

// type notifInfo struct {
// 	SenderID   int
// 	ReceiverID int
// 	Type       string
// 	Status     string
// 	Date       string
// }

// type notifData struct {
// 	Sender string
// 	SenderID int
// 	Avatar *string
// 	Type string
// 	Status string
// 	Date string
// }

// func GetNotifications(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("GetNotifications called")
// 	var logged_user loggedUser

// 	err := json.NewDecoder(r.Body).Decode(&logged_user)
// 	if err != nil {
// 		fmt.Println("Invalid Json:", err)
// 		return
// 	}

// 	var userId int
// 	err = dataB.SocialDB.QueryRow(`SELECT userId FROM Sessions WHERE id=?`, logged_user.UserSession).Scan(&userId)
// 	if err != nil {
// 		fmt.Println("Error get ID:", err)
// 		return
// 	}

// 	var notifications []interface{}
// 	rows, errf := dataB.SocialDB.Query(`SELECT senderId, receiverId, type, status, notificationDate FROM Notifications WHERE receiverId=?`, userId)
// 	if errf != nil {
// 		fmt.Println("Error fetching notifications:", errf)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var notifInfo notifInfo
// 		err = rows.Scan(&notifInfo.SenderID, &notifInfo.ReceiverID, &notifInfo.Type, &notifInfo.Status, &notifInfo.Date)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		var notifData notifData
// 		var firstname, lastname string
// 		var avatar *string
// 		err = dataB.SocialDB.QueryRow(`SELECT firstName, lastName, avatar FROM Users WHERE id=?`, notifInfo.SenderID).Scan(&firstname, &lastname, &avatar)
// 		if err != nil {
// 			fmt.Println("Error fetching user data:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		notifData.Sender = firstname+" "+lastname
// 		notifData.Type = notifInfo.Type
// 		notifData.Status = notifInfo.Status
// 		notifData.Date = notifInfo.Date
// 		notifData.SenderID = notifInfo.SenderID
// 		notifData.Avatar = avatar

// 		notifications = append(notifications, notifData)
// 	}

// 	response := map[string]interface{}{
// 		"notif_data":   notifications,
// 	}
// 	json.NewEncoder(w).Encode(response)
// }
