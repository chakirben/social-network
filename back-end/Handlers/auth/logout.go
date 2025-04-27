package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "socialN/dataBase"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to read cookie", http.StatusInternalServerError)
		return
	}

	sessionID := cookie.Value

	_, err = db.SocialDB.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		log.Println("Failed to delete session:", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	fmt.Fprintln(w, "Logged out successfully")
}
