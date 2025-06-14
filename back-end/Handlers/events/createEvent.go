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
		SELECT e.id, e.title, e.description, e.eventDate,
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

	rowss, err := database.SocialDB.Query("SELECT memberId FROM GroupsMembers WHERE groupId = ?", newEvent.GroupId)
	if err != nil {
		fmt.Println("Error fetching group members:", err)
	}
	fmt.Println("Group members fetched successfully")
	for rowss.Next() {
    var memberId int
    if err := rowss.Scan(&memberId); err != nil {
        fmt.Println("Error scanning memberId:", err)
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
                Type:             "Notification",
                NotificationType: "event",
                Title:            createdEvent.Title,
                Description:      createdEvent.Description,
                EventDate:        createdEvent.EventDate.String(),
                FirstName:        createdEvent.FirstName,
                LastName:         createdEvent.LastName,
                Avatar:           avatarStr,
            })
        }
    }
}


	if err := rowss.Err(); err != nil {
		fmt.Println("Error reading group members:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

type EventNotification struct {
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

func SendMessage(conn *websocket.Conn, msg EventNotification) {
	fmt.Println("Sending message to client:", msg)
	message, err := json.Marshal(msg)
	if err != nil {
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}
