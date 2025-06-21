package events

import (
	"encoding/json"
	"net/http"
	"time"

	"socialN/app/auth"
	database "socialN/dataBase"
)

type GetEvents struct {
	Id           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	EventDate    *time.Time `json:"eventDate"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Avatar       *string    `json:"avatar"`
	CreatorId    int        `json:"creatorId"`
	GroupId      int        `json:"groupId"`
	GoingMembers int        `json:"goingMembers"`
	IsUserGoing  *bool      `json:"isUserGoing"`
}

func GetHomeEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT 
		e.id, e.title, e.description, e.eventDate, e.creatorId,
		u.firstName, u.lastName, u.avatar,
		e.groupId,
		(SELECT COUNT(*) FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.isGoing = true) AS goingMembers,
		ea.isGoing AS isUserGoing
	FROM Events e
	JOIN Users u ON e.creatorId = u.id
	JOIN GroupsMembers gm ON gm.groupId = e.groupId AND gm.memberId = ?
	JOIN EventsAttendance ea ON ea.eventId = e.id AND ea.memberId = ? AND ea.isGoing = true
	ORDER BY e.eventDate ASC
	`

	rows, err := database.SocialDB.Query(query, userID, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve events: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []GetEvents

	for rows.Next() {
		var e GetEvents
		err := rows.Scan(
			&e.Id, &e.Title, &e.Description, &e.EventDate, &e.CreatorId,
			&e.FirstName, &e.LastName, &e.Avatar, &e.GroupId,
			&e.GoingMembers, &e.IsUserGoing,
		)
		if err != nil {
			http.Error(w, "Error scanning event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
