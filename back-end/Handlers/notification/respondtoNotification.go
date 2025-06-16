package notification

import (
	"net/http"
	"strconv"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

func RespondtoNotification(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	notificationIdStr := r.URL.Query().Get("notificationId")
	actionType := r.URL.Query().Get("actionType")

	notificationId, err := strconv.Atoi(notificationIdStr)
	if err != nil || notificationId < 1 {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if actionType != "accept" && actionType != "decline" {
		http.Error(w, "Invalid action type", http.StatusBadRequest)
		return
	}

	var notifType string
	var senderId, receiverId, groupId, eventId int
	err = dataB.SocialDB.QueryRow(`
		SELECT type, senderId, receiverId, groupId, eventId
		FROM Notifications
		WHERE id = ? AND status = 'pending'`,
		notificationId,
	).Scan(&notifType, &senderId, &receiverId, &groupId, &eventId)

	if err != nil {
		http.Error(w, "Notification not found or already handled", http.StatusNotFound)
		return
	}

	if userId != receiverId {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	if actionType == "accept" {
		switch notifType {
		case "follow_request":
			_, err = dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO Followers (followerId, followedId)
				VALUES (?, ?)`, senderId, receiverId)

		case "new_event":
			_, err = dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO EventsAttendance (memberId, eventId, isGoing)
				VALUES (?, ?, ?)`, receiverId, eventId, true)

		case "group_invite":
			_, err = dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO GroupsMembers (memberId, groupId)
				VALUES (?, ?)`, receiverId, groupId)

		case "group_join_request":
			_, err = dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO GroupsMembers (memberId, groupId)
				VALUES (?, ?)`, senderId, groupId)

		default:
			http.Error(w, "Unknown notification type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Failed to process acceptance", http.StatusInternalServerError)
			return
		}
	}

	_, err = dataB.SocialDB.Exec(`DELETE FROM Notifications WHERE id = ?`, notificationId)
	if err != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification handled successfully"))
}