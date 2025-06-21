package users

import (
	"encoding/json"
	"log"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type UnfollowedUser struct {
	ID            int     `json:"id"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	About         *string `json:"about"`
	Avatar        string  `json:"avatar"`
	FollowerCount int     `json:"followerCount"`
	HasRequested  bool    `json:"hasRequested"`
}

func GetUnfollowedUsers(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT u.id, u.firstName, u.lastName, u.about, u.avatar, 
	   EXISTS (
		SELECT 1 FROM Notifications n 
		WHERE n.receiverId = u.id AND n.senderId = ?
	) AS hasRequested,
	   COUNT(f2.followerId) as followerCount
	FROM Users u
	LEFT JOIN Followers f2 ON f2.followedId = u.id
	WHERE u.id != ? AND u.id NOT IN (
		SELECT followedId FROM Followers WHERE followerId = ?
	)
	GROUP BY u.id, u.firstName, u.lastName, u.about, u.avatar
	`

	// Pass all 3 parameters in the correct order:
	rows, err := dataB.SocialDB.Query(query, userID, userID, userID)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var unfollowedUsers []UnfollowedUser
	for rows.Next() {
		var user UnfollowedUser
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.About,
			&user.Avatar,
			&user.HasRequested,
			&user.FollowerCount,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		unfollowedUsers = append(unfollowedUsers, user)
	}

	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(unfollowedUsers)
	if err != nil {
		log.Println("JSON encoding error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
