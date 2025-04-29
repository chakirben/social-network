package groups

import (
	"encoding/json"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Insert the groups in the database....
func Creat_Groups(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var group Group

	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if group.Title == "" || group.Description == "" {
		http.Error(w, "Invalid JSON Title and Description are required... :(", http.StatusBadRequest)
		return
	}

	query := `
	  INSERT INTO groups (title , description , adminId) VALUES (?,?,?);
	`

	_, err = dataB.SocialDB.Exec(query, group.Title, group.Description, userID)
	if err != nil {
		log.Println("Error to insert groups in db :(", err)
		http.Error(w, "Failed to create group. Please try again later. :(", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
