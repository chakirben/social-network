package groups

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/app/auth"
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
            FROM Users u
            WHERE (
                EXISTS (
                    SELECT 1 FROM Followers f
                    WHERE (f.followerId = ? AND f.followedId = u.id)
                       OR (f.followerId = u.id AND f.followedId = ?)
                )
            )
            AND NOT EXISTS (
                SELECT 1 FROM GroupsMembers gm
                WHERE gm.memberId = u.id AND gm.groupId = ?
            )    
	`

	rows, err := dataB.SocialDB.Query(query, userID, userID, Group.GroupID)
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
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var status sql.NullString
		checkQuery := `SELECT status FROM Notifications WHERE groupId = ? AND type = 'group_invite' AND senderId = ? AND receiverId = ? LIMIT 1`
		err := dataB.SocialDB.QueryRow(checkQuery, Group.GroupID, userID, u.ID).Scan(&status)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("error checking notification for group:", err)
			http.Error(w, "Error checking notification", http.StatusInternalServerError)
			return
		}

		if status.Valid {
			u.Status = "Cancel-Invite"
		} else {
			u.Status = "+invite"
		}

		users = append(users, u)
	}

	fmt.Println("gg", users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
