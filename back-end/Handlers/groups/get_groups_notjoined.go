package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataB "socialN/dataBase"
)

type NotMyGroups struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	Status bool
}

// Get all groups that the user has not joined yet...
func GetGroupsUserNotJoined(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	// 	return
	// }
	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Invalid session :(", http.StatusUnauthorized)
	// 	return
	// }

	query := `
	SELECT 
		g.id,
		g.title,
		g.description,
		COUNT(gm.memberId) AS members_count
	FROM Groups g
	LEFT JOIN GroupsMembers gm ON g.id = gm.groupId
	WHERE g.id NOT IN (
		SELECT groupId 
		FROM GroupsMembers 
		WHERE memberId = ?
	)
	GROUP BY g.id, g.title, g.description
	`

	rows, err := dataB.SocialDB.Query(query, 1)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var groups []NotMyGroups
	for rows.Next() {
		var g NotMyGroups
		if err := rows.Scan(&g.Id, &g.Title, &g.Description, &g.MembersCount); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}
		g.Status = false
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
