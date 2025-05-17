package events

import (
	"encoding/json"
	"net/http"
	"socialN/Handlers/auth"
	database "socialN/dataBase"
)

func GetUserEventsHandler(w http.ResponseWriter, r *http.Request) {
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
	JOIN GroupsMembers gm ON e.groupId = gm.groupId
	JOIN Users u ON u.id = e.creatorId
	WHERE gm.memberId = ?
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
