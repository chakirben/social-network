package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	requesterID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	otherUserIdStr := r.URL.Query().Get("id")
	var otherUserID int
	if otherUserIdStr == "" {
		otherUserID = requesterID
	} else {
		otherUserID, err = strconv.Atoi(otherUserIdStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	}

	listType := r.URL.Query().Get("type")
	if listType != "followers" && listType != "following" {
		http.Error(w, "Invalid list type", http.StatusBadRequest)
		return
	}

	var accountType string
	err = dataB.SocialDB.QueryRow("SELECT accountType FROM Users WHERE id = ?", otherUserID).Scan(&accountType)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if accountType == "private" && requesterID != otherUserID {
		var isFollower int
		err = dataB.SocialDB.QueryRow(
			"SELECT COUNT(*) FROM Followers WHERE followerId = ? AND followedId = ?",
			requesterID, otherUserID,
		).Scan(&isFollower)
		if err != nil {
			log.Printf("Error querying follower status: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if isFollower == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Access denied. You need to follow this user to view their followers/following."})
			return
		}
	}

	var query string
	if listType == "followers" {
		query = `
			SELECT u.id, u.firstName, u.lastName, u.avatar
			FROM Users u
			INNER JOIN Followers f ON u.id = f.followerId
			WHERE f.followedId = ?`
	} else {
		query = `
			SELECT u.id, u.firstName, u.lastName, u.avatar
			FROM Users u
			INNER JOIN Followers f ON u.id = f.followedId
			WHERE f.followerId = ?`
	}

	rows, err := dataB.SocialDB.Query(query, otherUserID)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName string
		var avatar *string
		if err := rows.Scan(&id, &firstName, &lastName, &avatar); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		users = append(users, map[string]interface{}{
			"id":        id,
			"firstName": firstName,
			"lastName":  lastName,
			"avatar":    avatar,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
