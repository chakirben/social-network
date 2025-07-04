package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type MyGroups struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	PostCont     int
}

// Get all groups that the user has joined...
func GetMyGroups(w http.ResponseWriter, r *http.Request) {
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
		COUNT(gm2.memberId) AS members_count,
		COUNT(P.groupId) AS CP
	FROM Groups g
	JOIN GroupsMembers gm1 ON g.id = gm1.groupId
	LEFT JOIN GroupsMembers gm2 ON g.id = gm2.groupId
	LEFT JOIN Posts P ON P.groupId = g.id
	WHERE gm1.memberId = ?
	GROUP BY g.id, g.title, g.description
	`

	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var groups []MyGroups
	for rows.Next() {
		var g MyGroups
		if err := rows.Scan(&g.Id, &g.Title, &g.Description, &g.MembersCount, &g.PostCont); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
