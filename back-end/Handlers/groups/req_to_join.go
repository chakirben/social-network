package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	SendToAdmin := struct {
		Admin  int
		GroupID int
		Userid int
	}{
		Admin:  adminId,
		GroupID:  req.GroupID,
		Userid: userID,
	}
	// this is a just a test until we  creat ws ...
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "applicaton/json")
	if err := json.NewEncoder(w).Encode(SendToAdmin); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
