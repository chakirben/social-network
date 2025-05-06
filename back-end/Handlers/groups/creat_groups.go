package groups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dataB "socialN/dataBase"
)

type Grouppp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Insert the groups in the database....
func Creat_Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Invalid session :(", http.StatusUnauthorized)
	// 	return
	// }

	var group Grouppp

	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		fmt.Println("hiiii yosf")
		
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

	ress, err := dataB.SocialDB.Exec(query, group.Title, group.Description, 2)
	if err != nil {
		log.Println("Error to insert groups in db :(", err)
		http.Error(w, "Failed to create group. Please try again later. :(", http.StatusInternalServerError)
		return
	}

	lastIDgroup, err := ress.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to create group. Please try again later. :(", http.StatusInternalServerError)
		log.Fatal(err)
	}

	queryy := `
		INSERT INTO GroupsMembers (memberId, groupId) VALUES (?, ?)
	`
	_, err = dataB.SocialDB.Exec(queryy, 1, lastIDgroup)
	if err != nil {
		log.Println("Error to insert members in db :(", err)
		http.Error(w, "Failed to create group. Please try again later. :(", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
