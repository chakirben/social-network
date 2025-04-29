package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dataB "socialN/dataBase"
)

type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func Creat_Groups(w http.ResponseWriter, r *http.Request) {
	fmt.Println("kdhlt")
	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Invalid session", http.StatusUnauthorized)
	// 	return
	// }

	var group Group

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		 http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	query := `
	  INSERT INTO groups (title , description , adminId) VALUES (?,?,?);
	`

	_, err = dataB.SocialDB.Exec(query, group.Title, group.Description, 1)
	if err != nil {
		log.Println("Error to insert groups in db:", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
