package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "socialN/dataBase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LogUser(w http.ResponseWriter, r *http.Request) {
	db := db.SocialDB
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	var userID int
	var passwordHash string
	query := `SELECT id, password FROM Users WHERE email = ?`
	err := db.QueryRow(query, loginRequest.Email).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessionID, expiration, err := CreateSession(db, userID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful")
}

func CreateSession(db *sql.DB, userID int) (string, time.Time, error) {
	sessionID := uuid.New().String()
	expiration := time.Now().AddDate(1000, 0, 0)
	query := `INSERT INTO Sessions (id, userId, expiresAt) VALUES (?, ?, ?)`
	_, err := db.Exec(query, sessionID, userID, expiration)
	if err != nil {
		log.Println("Error storing session in database:", err)
		return "", time.Time{}, err
	}
	return sessionID, expiration, nil
}
