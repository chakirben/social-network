package events

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sess "socialN/Handlers/auth"
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

	userID, err := sess.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "err"+err.Error(), http.StatusUnauthorized)
		return
	}

	var newEvent Event
	newEvent.CreatorId = userID
	err = json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var grID int
	err = database.SocialDB.QueryRow(`SELECT 1 FROM GroupsMembers WHERE memberId = ? AND groupId = ? LIMIT 1`, userID, newEvent.GroupId).Scan(&grID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not a member of this group", http.StatusForbidden)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("Error during query execution:", err)
		return
	}

	newEvent.CreatedAt = time.Now()

	if newEvent.Title == "" || newEvent.Description == "" {
		http.Error(w, "missing required Field", http.StatusBadRequest)
		return
	}

	result, err := database.SocialDB.Exec(`INSERT INTO Events (title, description, eventDate, creatorId, groupId, createdAt) VALUES (?, ?, ?, ?, ?, ?)`,
		newEvent.Title, newEvent.Description, newEvent.EventDate, newEvent.CreatorId, newEvent.GroupId, newEvent.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to insert event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve the event ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now, retrieve the necessary data to return in the response
	var creatorFirstName, creatorLastName, creatorAvatar string
	var eventTitle, eventDescription string
	var eventDate time.Time

	err = database.SocialDB.QueryRow(`
		SELECT u.firstName, u.lastName, u.avatar, e.title, e.description, e.eventDate
		FROM Events e
		JOIN Users u ON e.creatorId = u.id
		WHERE e.id = ?`, eventID).Scan(&creatorFirstName, &creatorLastName, &creatorAvatar, &eventTitle, &eventDescription, &eventDate)
	if err != nil {
		http.Error(w, "Failed to fetch event data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(creatorFirstName)
	response := map[string]interface{}{
		"title":            eventTitle,
		"description":      eventDescription,
		"eventDate":        eventDate,
		"creatorFirstName": creatorFirstName,
		"creatorLastName":  creatorLastName,
		"creatorAvatar":    creatorAvatar,
	}
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// func InsertEvent(e Event) error {
// 	_, err := database.SocialDB.Exec(`INSERT INTO Events (title, description, eventDate, creatorId, groupId, createdAt) VALUES (?, ?, ?, ?, ?, ?)`,
// 		e.Title, e.Description, e.EventDate, e.CreatorId, e.GroupId, e.CreatedAt)
// 	return err
// }
