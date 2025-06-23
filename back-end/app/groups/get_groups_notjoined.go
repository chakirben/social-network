package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type NotMyGroups struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	Status       string
	PostCont     int
}

// Get all groups that the user has not joined yet...
func GetGroupsUserNotJoined(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT 
		g.id,
		g.title,
		g.description,
		COUNT(gm.memberId) AS members_count,
		COUNT(P.groupId) AS CP
	FROM Groups g
	LEFT JOIN GroupsMembers gm ON g.id = gm.groupId
	LEFT JOIN Posts P ON P.groupId = g.id
	WHERE g.id NOT IN (
		SELECT groupId 
		FROM GroupsMembers 
		WHERE memberId = ?
	)
	GROUP BY g.id, g.title, g.description
	`

	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var groups []NotMyGroups
	for rows.Next() {
		var g NotMyGroups
		if err := rows.Scan(&g.Id, &g.Title, &g.Description, &g.MembersCount, &g.PostCont); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM Notifications WHERE senderId = ? AND groupId = ?  AND  type = "group_join_request"  LIMIT 1);`
		err := dataB.SocialDB.QueryRow(checkQuery, userID, g.Id).Scan(&exists)
		if err != nil {
			fmt.Println("error checking notification for group:", err)
			http.Error(w, "error checking notification", http.StatusInternalServerError)
			return
		}

		if exists {
			g.Status = "Cancel"
		} else {
			g.Status = "Join"
		}

		groups = append(groups, g)

	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
