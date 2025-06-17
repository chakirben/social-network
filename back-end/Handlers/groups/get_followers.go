package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Followers struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
}


type GroupInvite struct {
	GroupID int `json:"groupId"`
}

func GetFollowersList(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var Group GroupInvite
	err = json.NewDecoder(r.Body).Decode(&Group)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}


	query := `
		SELECT u.id, u.firstName, u.lastName, u.avatar
		FROM Followers f
		JOIN Users u ON f.followedId = u.id
		WHERE f.followerId = ?
	`
	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		log.Println("DB query error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []Followers
	for rows.Next() {
		var u Followers
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar); err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM Notifications WHERE groupId = ?  AND  type = "group_invite" AND senderId = ? AND receiverId = ? LIMIT 1);`
		err := dataB.SocialDB.QueryRow(checkQuery, Group.GroupID , userID , u.ID).Scan(&exists)
		if err != nil {
			fmt.Println("error checking notification for group:", err)
			http.Error(w, "error checking notification", http.StatusInternalServerError)
			return
		}

		if exists {
			u.Status = "INVITE"
		} else {
			u.Status = "CancelInvite"
		}

		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
