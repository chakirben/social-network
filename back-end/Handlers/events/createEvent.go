package events

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sess "socialN/Handlers/auth"
	"socialN/Handlers/ws"
	database "socialN/dataBase"

	"github.com/gorilla/websocket"
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

type EventInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"eventDate"`
	GroupId     int    `json:"groupId"`
}

type EventNotification struct {
	Id               int64  `json:"id"`
	Type             string `json:"type"`
	NotificationType string `json:"notificationType"`
	ItemId           int    `json:"itemId"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	EventDate        string `json:"eventDate"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Avatar           string `json:"avatar"`
}

type CreatedEvent struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	EventDate    time.Time `json:"eventDate"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Avatar       *string   `json:"avatar"`
	GroupId      int       `json:"groupId"`
	GoingMembers int       `json:"goingMembers"`
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
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.Title == "" || input.Description == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		http.Error(w, "Invalid date format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if parsedDate.Before(time.Now()) {
		http.Error(w, "Event date must be in the future", http.StatusBadRequest)
		return
	}

	// Check if the user is in the group
	var exists int
	err = database.SocialDB.QueryRow(
		`SELECT 1 FROM GroupsMembers WHERE memberId = ? AND groupId = ? LIMIT 1`,
		userID, input.GroupId,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not a member of this group", http.StatusForbidden)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Insert event
	result, err := database.SocialDB.Exec(
		`INSERT INTO Events (title, description, eventDate, creatorId, groupId, createdAt)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		input.Title, input.Description, parsedDate, userID, input.GroupId, time.Now(),
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

	// Create notifications for group members (excluding creator)
	notifResult, err := database.SocialDB.Exec(`
		INSERT INTO Notifications (type, senderId, receiverId, groupId, eventId, status)
		SELECT 'new_event', ?, memberId, ?, ?, 'pending'
		FROM GroupsMembers
		WHERE groupId = ? AND memberId != ?`,
		userID, input.GroupId, eventID, input.GroupId, userID,
	)
	if err != nil {
		http.Error(w, "Failed to insert notifications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	notificationID, err := notifResult.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve notification ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch event details for WebSocket
	var createdEvent CreatedEvent
	query := `
		SELECT e.id, e.title, e.description, e.eventDate,
		       u.firstName, u.lastName, u.avatar, e.groupId,
		       (SELECT COUNT(*) FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.isGoing = true) AS goingMembers
		FROM Events e
		JOIN Users u ON u.id = e.creatorId
		WHERE e.id = ?
	`
	err = database.SocialDB.QueryRow(query, eventID).Scan(
		&createdEvent.Id,
		&createdEvent.Title,
		&createdEvent.Description,
		&createdEvent.EventDate,
		&createdEvent.FirstName,
		&createdEvent.LastName,
		&createdEvent.Avatar,
		&createdEvent.GroupId,
		&createdEvent.GoingMembers,
	)
	if err != nil {
		http.Error(w, "Failed to fetch created event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send WebSocket notifications to all group members
	rows, err := database.SocialDB.Query("SELECT memberId FROM GroupsMembers WHERE groupId = ?", input.GroupId)
	if err != nil {
		http.Error(w, "Failed to fetch group members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var memberId int
		if err := rows.Scan(&memberId); err != nil || memberId == userID {
			continue
		}

		ws.ConnMu.Lock()
		conns := ws.Connections[memberId]
		ws.ConnMu.Unlock()

		for _, conn := range conns {
			if conn != nil {
				avatarStr := ""
				if createdEvent.Avatar != nil {
					avatarStr = *createdEvent.Avatar
				}

				SendMessage(conn, EventNotification{
					Id:               notificationID,
					Type:             "Notification",
					NotificationType: "event",
					ItemId:           createdEvent.Id,
					Title:            createdEvent.Title,
					Description:      createdEvent.Description,
					EventDate:        createdEvent.EventDate.Format(time.RFC3339),
					FirstName:        createdEvent.FirstName,
					LastName:         createdEvent.LastName,
					Avatar:           avatarStr,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

func SendMessage(conn *websocket.Conn, msg EventNotification) {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to marshal message:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		fmt.Println("Error sending WebSocket message:", err)
	}
}
