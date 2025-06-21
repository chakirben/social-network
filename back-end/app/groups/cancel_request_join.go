package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type CancelReqJoin struct {
	GroupID int `json:"groupId"`
}

func CancelRequestToJoinGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var req CancelReqJoin
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
       DELETE FROM Notifications WHERE senderId = ? AND type = "group_join_request" AND  groupId = ?
    `

	_, err = dataB.SocialDB.Exec(query2, userID, req.GroupID)
	if err != nil {
		fmt.Println(err)
		log.Println("Error to select admin in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
