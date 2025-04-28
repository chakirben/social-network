package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "socialN/dataBase"
)

type UserProfile struct {
	Nickname  string
	Age       int
	FirstName string
	LastName  string
	Email     string
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := ValidateSession(r, db.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var uu UserProfile
	errD := db.SocialDB.QueryRow(`SELECT nickname, email , age, firstName, lastName FROM Users WHERE id = ?`,
		userID).Scan(&uu.Nickname, &uu.Email, &uu.Age, &uu.FirstName, &uu.LastName)
	if errD == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if errD != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(uu)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
