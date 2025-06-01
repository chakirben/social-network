package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type ReqJoinGroup struct {
	GroupID int `json:"groupId"`
}

func Req_To_Join_Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var req ReqJoinGroup
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if req.GroupID < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}

	var adminId int
	query := `SELECT adminId FROM Groups WHERE id = ?`
	err = dataB.SocialDB.QueryRow(query, req.GroupID).Scan(&adminId)
	if err != nil {
		fmt.Println(err)
		log.Println("Error to select admin in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	query2 := `
	   INSERT INTO Notifications (senderId, receiverId, type, notificationDate , groupTargetId) VALUES (?,?,?,?,?)
	`

	_, err = dataB.SocialDB.Exec(query2, userID, adminId, "group_join_request", time.Now(), req.GroupID)
	if err != nil {
		fmt.Println(err)
		log.Println("Error to select admin in db :(", err)
		http.Error(w, "Failed to join group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "applicaton/json")
	if err := json.NewEncoder(w).Encode("pending"); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
