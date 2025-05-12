package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type user struct {
	ID          int            `json:"id"`
	Nickname    *string        `json:"nickname"`
	Email       string         `json:"email"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	DateOfBirth string         `json:"dateOfBirth"`
	Avatar      *string        `json:"avatar"`
	About       *string        `json:"about"`
	AccountType sql.NullString `json:"accountType"`
}

func GetCurrentUserData(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var u user
	err = dataB.SocialDB.QueryRow(
		`SELECT id, nickname, email, firstName, lastName, dateOfBirth, avatar, about, accountType
		FROM Users WHERE id = ?`, userID).Scan(
		&u.ID, &u.Nickname, &u.Email, &u.FirstName, &u.LastName,
		&u.DateOfBirth, &u.Avatar, &u.About, &u.AccountType,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
