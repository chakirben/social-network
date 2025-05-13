package users

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

func Follow(w http.ResponseWriter, r *http.Request) {
	senderId, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverIdStr := r.URL.Query().Get("id")
	if receiverIdStr == "" {
		http.Error(w, "Missing receiver id", http.StatusBadRequest)
		return
	}

	receiverId, err := strconv.Atoi(receiverIdStr)
	if err != nil {
		http.Error(w, "Invalid receiver id", http.StatusBadRequest)
		return
	}

	if senderId == receiverId {
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if a follow request already exists
	var exists int
	err = dataB.SocialDB.QueryRow(`
		SELECT 1 FROM Notifications 
		WHERE senderId = ? AND receiverId = ? AND type = 'follow_request'
	`, senderId, receiverId).Scan(&exists)

	if err == nil {
		// Exists → Unfollow (delete the row)
		_, err = dataB.SocialDB.Exec(`
			DELETE FROM Notifications 
			WHERE senderId = ? AND receiverId = ? AND type = 'follow_request'
		`, senderId, receiverId)
		if err != nil {
			log.Println("Error deleting notification:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unfollowed successfully"))
		return
	}

	// Does not exist → Follow (insert the row)
	_, err = dataB.SocialDB.Exec(`
		INSERT INTO Notifications (senderId, receiverId, type, status, notificationDate)
		VALUES (?, ?, 'follow_request', 'pending', ?)
	`, senderId, receiverId, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Error inserting notification:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Follow request sent successfully"))
}
