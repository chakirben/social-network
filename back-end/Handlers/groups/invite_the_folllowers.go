package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type TheUserInvited struct {
	UserId  int `json:"userId"`  
	GroupId int `json:"groupId"`  
}

func InfiteTheFollowers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var userIDToInvite TheUserInvited
	err = json.NewDecoder(r.Body).Decode(&userIDToInvite)
	if err != nil {
		fmt.Println("GGHHH hi ysosf", err)
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if userIDToInvite.UserId < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}
	query := `
	INSERT INTO Notifications (senderId, receiverId , type , groupId )
	VALUES (?, ?, ?, ?)
`

	_, err = dataB.SocialDB.Exec(query, userID, userIDToInvite.UserId, "group_invite", userIDToInvite.GroupId)
	if err != nil {
		http.Error(w, "Could not send invitation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
