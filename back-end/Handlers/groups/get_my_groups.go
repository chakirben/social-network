package groups

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Groups struct {
	Id          int
	Title       string
	Description string
	Members     int
}

func GetMyGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT id , title , description 
	FROM Groups g 
	JOIN GroupsMembers gm ON g.id = gm.groupId
	WHERE gm.memberId = ?
	`

	rows, err := dataB.SocialDB.Query(query, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var groups []Groups
	for rows.Next() {
		var g Groups
		if err := rows.Scan(&g.Id, &g.Title, &g.Description); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}
		members, err := GetMembersGroups(g.Id)
		if err != nil {
			fmt.Println("error to get members", err)
			http.Error(w, "error to get members of groups", http.StatusInternalServerError)
			return
		}
		g.Members = members
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
