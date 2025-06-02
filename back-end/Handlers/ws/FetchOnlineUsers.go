package ws

import (
	"encoding/json"
	"net/http"
	"strings"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	ids := []int{}
	for id := range Connections {
		if id != userID { 
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1] 

	query := "SELECT id, firstName, lastName, avatar FROM Users WHERE id IN (" + placeholders + ")"

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := dataB.SocialDB.Query(query, args...)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Avatar    string `json:"avatar"`
		Status    string `json:"status"`
	}

	var onlineUsers []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Avatar)
		if err != nil {
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			return
		}
		u.Status = "online"
		onlineUsers = append(onlineUsers, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(onlineUsers)
}
