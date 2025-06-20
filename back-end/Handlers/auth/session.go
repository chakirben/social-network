package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	dataB "socialN/dataBase"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	id, err := ValidateSession(r, dataB.SocialDB)
	fmt.Println(id, err)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(id)
	}

func ValidateSession(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("ysf", err)
			return 0, fmt.Errorf("no session cookie")
		}
		return 0, fmt.Errorf("error reading session cookie: %v", err)
	}

	sessionID := cookie.Value
	var userID int
	var expiresAt time.Time
	query := `SELECT userId, expiresAt FROM Sessions WHERE id = ?`
	err = db.QueryRow(query, sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, fmt.Errorf("session not found: %v", err)
	}
	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("session expired")
	}
	
	return userID, nil
}
