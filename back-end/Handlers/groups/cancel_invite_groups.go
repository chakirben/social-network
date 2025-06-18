package groups

import (
	"encoding/json"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type CancelInvite struct {
	GroupID    int `json:"groupId"`
	ReceiverId int `json:"userid"`
}

func CancelInviteToGroups(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var req CancelInvite
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if req.GroupID < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	query2 := `
       DELETE FROM Notifications WHERE senderId = ? AND receiverId = ? AND groupId = ? AND type = 'group_invite'
    `

	_, err = dataB.SocialDB.Exec(query2, userID, req.ReceiverId, req.GroupID)
	if err != nil {
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
