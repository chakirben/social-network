package events

import (
	"database/sql"
	"encoding/json"
	"net/http"

	database "socialN/dataBase"
)

var eveAttendence struct {
	EventId int  `json:"eventId"`
	GroupId int  `json:"group_id"`
	IsGoing bool `json:"isGoing"`
}

func SetAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	// userID, err := auth.ValidateSession(r, database.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	userID := 1
	if err := json.NewDecoder(r.Body).Decode(&eveAttendence); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var grID int
	err := database.SocialDB.QueryRow(`SELECT 1 FROM GroupsMembers WHERE memberId = ? AND groupId = ? LIMIT 1`, userID, eveAttendence.GroupId).Scan(&grID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not a member of this group", http.StatusForbidden)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = database.SocialDB.Exec(`INSERT INTO EventsAttendance(memberId, eventId, IsGoing) VALUES (?, ?, ?)
									  ON CONFLICT(memberId, eventId) DO UPDATE SET isGoing = excluded.isGoing;`,
		userID, eveAttendence.EventId, eveAttendence.IsGoing)
	if err != nil {
		http.Error(w, "db Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Attendance recorded"))
}
