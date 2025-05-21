package events

import (
	"encoding/json"
	"net/http"

	"socialN/Handlers/auth"
	database "socialN/dataBase"
)

func GetGroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT e.id, e.title, e.description, e.eventDate, e.creatorId,
	       u.firstName, u.lastName, u.avatar, e.groupId,
	       (SELECT COUNT(*) FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.isGoing = true) AS goingMembers,
	       (SELECT ea.isGoing FROM EventsAttendance ea WHERE ea.eventId = e.id AND ea.memberId = ? LIMIT 1) AS isUserGoing
	FROM Events e
	JOIN Users u ON u.id = e.creatorId
	WHERE e.groupId = ?
	ORDER BY e.eventDate DESC
	`

	rows, err := database.SocialDB.Query(query, userID, groupID)
	if err != nil {
		http.Error(w, "Failed to retrieve events: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []GetEvents
	for rows.Next() {
		var e GetEvents
		err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.EventDate, &e.CreatorId, &e.FirstName, &e.LastName, &e.Avatar, &e.GroupId, &e.GoingMembers, &e.IsUserGoing)
		if err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
