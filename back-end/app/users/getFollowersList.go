package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type Follower struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
}

func GetFollowersListHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	query := `
		SELECT u.id, u.firstName, u.lastName, u.avatar
		FROM Followers f
		JOIN Users u ON f.followerId = u.id
		WHERE f.followedId = ?
	`
	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		log.Println("DB query error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []Follower
	for rows.Next() {
		var u Follower
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar); err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	fmt.Println("Followers List:", users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
