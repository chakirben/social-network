package events

import (
	"database/sql"
	"encoding/json"
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

// A separate input struct to decode raw ISO date string
type EventInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"eventDate"`
	GroupId     int    `json:"groupId"`
}

func SetEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := sess.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var input EventInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.Title == "" || input.Description == "" {
		http.Error(w, "Missing required field", http.StatusBadRequest)
		return
	}

	// Parse ISO 8601 date from frontend
	parsedDate, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		http.Error(w, "Invalid date format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if parsedDate.Before(time.Now()) {
		http.Error(w, "Event date must be in the future", http.StatusBadRequest)
		return
	}

	newEvent := Event{
		Title:       input.Title,
		Description: input.Description,
		EventDate:   parsedDate,
		GroupId:     input.GroupId,
		CreatorId:   userID,
		CreatedAt:   time.Now(),
	}

	var exists int
	err = database.SocialDB.QueryRow(
		`SELECT 1 FROM GroupsMembers WHERE memberId = ? AND groupId = ? LIMIT 1`,
		userID, newEvent.GroupId,
	).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not a member of this group", http.StatusForbidden)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	result, err := database.SocialDB.Exec(
		`INSERT INTO Events (title, description, eventDate, creatorId, groupId, createdAt)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		newEvent.Title, newEvent.Description, newEvent.EventDate, newEvent.CreatorId, newEvent.GroupId, newEvent.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to insert event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve event ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var createdEvent GetEvents
	query := `
		SELECT e.id, e.title, e.description, e.eventDate, e.creatorId,
		       u.firstName, u.lastName, u.avatar, e.groupId,
		       (SELECT COUNT(*) FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.isGoing = true) AS goingMembers,
		       IFNULL((SELECT ea.isGoing FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.memberId = ? LIMIT 1), false) AS isUserGoing
		FROM Events e
		JOIN Users u ON u.id = e.creatorId
		WHERE e.id = ?
	`

	err = database.SocialDB.QueryRow(query, userID, eventID).Scan(
		&createdEvent.Id,
		&createdEvent.Title,
		&createdEvent.Description,
		&createdEvent.EventDate,
		&createdEvent.CreatorId,
		&createdEvent.FirstName,
		&createdEvent.LastName,
		&createdEvent.Avatar,
		&createdEvent.GroupId,
		&createdEvent.GoingMembers,
		&createdEvent.IsUserGoing,
	)
	if err != nil {
		http.Error(w, "Failed to fetch created event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}
