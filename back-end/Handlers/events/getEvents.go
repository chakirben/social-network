package events

import (
	"encoding/json"
	"net/http"
	"time"

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

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID := 3
	query := `SELECT e.id, e.title, e.description, e.eventDate, e.creatorId,
					(SELECT u.firstName FROM Users u WHERE u.id = ?) AS firstName,
					(SELECT u.lastName FROM Users u WHERE u.id = ?) AS lastName, e.groupId,
					(SELECT u.avatar FROM Users u WHERE u.id = ?) AS avatar,
			(SELECT COUNT(*) FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.isGoing = true) AS goingMembers,
			(SELECT ea.isGoing FROM EventsAttendance ea WHERE ea.eventId = e.id AND  ea.memberId = ? LIMIT 1) AS isUserGoing
			FROM Events e`

	rows, err := database.SocialDB.Query(query, userID, userID, userID, userID)
	if err != nil {
		http.Error(w, "Failed retrieve events ---"+err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()
	var events []GetEvents

	for rows.Next() {
		var e GetEvents
		err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.EventDate, &e.CreatorId, &e.FirstName, &e.LastName, &e.GroupId, &e.Avatar, &e.GoingMembers, &e.IsUserGoing)
		if err != nil {
			http.Error(w, "Error scanning event-------: "+err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
	}
	// fmt.Println("eve", events)
}
