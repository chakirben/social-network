package events

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	// sess "socialN/Handlers/auth"
	database "socialN/dataBase"
)

type Event struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"eventDate"`
	CreatorId   int       `json:"creatorId"`
	GroupId     int       `json:"groupId"`
	CreatedAt   time.Time `json:"created_at"`
}

func SetEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// userID, err := sess.ValidateSession(r, database.SocialDB)
	// if err != nil {
	// 	http.Error(w, "err"+err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	userID := 1
	var newEvent Event

	newEvent.CreatorId = userID
	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var grID int
	err = database.SocialDB.QueryRow(`SELECT 1 FROM GroupsMembers WHERE memberId = ? AND groupId = ? LIMIT 1`, userID, newEvent.GroupId).Scan(&grID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not a member of this group", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newEvent.CreatedAt = time.Now()

	if newEvent.Title == "" || newEvent.Description == "" {
		http.Error(w, "missing required Field", http.StatusBadRequest)
		return
	}

	err = InsertEvent(newEvent)
	if err != nil {
		http.Error(w, "Failed to insert event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Event created successfully"))
}

func InsertEvent(e Event) error {
	_, err := database.SocialDB.Exec(`INSERT INTO Events (title, description, eventDate, creatorId, groupId, createdAt) VALUES (?, ?, ?, ?, ?, ?)`,
		e.Title, e.Description, e.EventDate, e.CreatorId, e.GroupId, e.CreatedAt)
	return err
}
