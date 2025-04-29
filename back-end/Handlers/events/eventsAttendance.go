package events

import (
	"encoding/json"
	"net/http"

	database "socialN/dataBase"
)

var eveAttendence struct {
	EventId int  `json:"eventId"`
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

	_, err := database.SocialDB.Exec(`INSERT INTO EventsAttendance(memberId, eventId, IsGoing) VALUES (?, ?, ?)
									  ON CONFLICT(memberId, eventId) DO UPDATE SET isGoing = excluded.isGoing;`,
		userID, eveAttendence.EventId, eveAttendence.IsGoing)
	if err != nil {
		http.Error(w, "db Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Attendance recorded"))
}
