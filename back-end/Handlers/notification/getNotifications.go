package notification

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataB "socialN/dataBase"
)

type loggedUser struct {
	UserSession string `json:"user_session"`
}

type notifInfo struct {
	SenderID   int
	ReceiverID int
	Type       string
	Status     string
	Date       string
}

type notifData struct {
	Sender string
	SenderID int
	Avatar *string
	Type string
	Status string
	Date string
}

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	var logged_user loggedUser

	err := json.NewDecoder(r.Body).Decode(&logged_user)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	//get user id
	var userId int
	err = dataB.SocialDB.QueryRow(`SELECT userId FROM Sessions WHERE id=?`, logged_user.UserSession).Scan(&userId)
	if err != nil {
		fmt.Println("Error get ID:", err)
		return
	}


	var notifications []interface{}
	//get Data from DB
	rows, errf := dataB.SocialDB.Query(`SELECT senderId, receiverId, type, status, notificationDate FROM Notifications WHERE receiverId=?`, userId)
	if errf != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var notifInfo notifInfo
		err = rows.Scan(&notifInfo.SenderID, &notifInfo.ReceiverID, &notifInfo.Type, &notifInfo.Status, &notifInfo.Date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var notifData notifData
		var firstname, lastname string
		var avatar *string
		err = dataB.SocialDB.QueryRow(`SELECT firstName, lastName, avatar FROM Users WHERE id=?`, notifInfo.SenderID).Scan(&firstname, &lastname, &avatar)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		notifData.Sender = firstname+" "+lastname
		notifData.Type = notifInfo.Type
		notifData.Status = notifInfo.Status
		notifData.Date = notifInfo.Date 
		notifData.SenderID = notifInfo.SenderID
		notifData.Avatar = avatar

		notifications = append(notifications, notifData)
	}


	//send response
	response := map[string]interface{}{
		"notif_data":   notifications,
	}
	json.NewEncoder(w).Encode(response)
}
