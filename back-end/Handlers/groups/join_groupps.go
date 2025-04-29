package groups

import (
	"encoding/json"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type JoinGroupRequest struct {
	GroupID int `json:"groupId"`
}

// Insert the GroupsMembers in the database....
func JoinGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var req JoinGroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if req.GroupID < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO GroupsMembers (memberId, groupId) VALUES (?, ?)
	`
	_, err = dataB.SocialDB.Exec(query, userID, req.GroupID)
	if err != nil {
		log.Println("Error to insert members in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
