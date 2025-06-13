package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Friend struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserConnectionsResponse struct {
	Friends []Friend `json:"friends"`
	Groups  []Group  `json:"groups"`
}

func GetUserConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println(1)
	friendRows, err := dataB.SocialDB.Query(`
		SELECT DISTINCT u.id, u.firstName, u.lastName, u.avatar
		FROM Users u
		JOIN Followers f ON 
			(f.followerId = u.id AND f.followedId = ?) OR 
			(f.followedId = u.id AND f.followerId = ?)
	`, userID, userID)
	if err != nil {
		log.Println("DB query error (friends):", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(2)

	defer friendRows.Close()

	var friends []Friend
	for friendRows.Next() {
		var f Friend
		if err := friendRows.Scan(&f.ID, &f.FirstName, &f.LastName, &f.Avatar); err != nil {
			log.Println("Friend scan error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		friends = append(friends, f)
	}
	groupRows, err := dataB.SocialDB.Query(`
		SELECT g.id, g.title
		FROM GroupsMembers gm
		JOIN Groups g ON gm.groupId = g.id
		WHERE gm.memberId = ?
	`, userID)
	if err != nil {
		log.Println("DB query error (groups):", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(3)
	defer groupRows.Close()
	var groups []Group
	for groupRows.Next() {
		var g Group
		if err := groupRows.Scan(&g.ID, &g.Name); err != nil {
			log.Println("Group scan error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		groups = append(groups, g)
	}
	response := UserConnectionsResponse{
		Friends: friends,
		Groups:  groups,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
