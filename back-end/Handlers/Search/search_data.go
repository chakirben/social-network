package searchdata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type NotfollowedUser struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	About         string `json:"about"`
	Avatar        string `json:"avatar"`
	FollowerCount int    `json:"followerCount"`
}

type NotJoinGroups struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	PostCont     int
}
type JoinedGroups struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	PostCont     int
}

type Data struct {
	Notfollowed  []NotfollowedUser `json:"Notfollowed"`
	UnJoinGroups []NotJoinGroups `json:"UnJoinGroups"`
	JoinedGroups []JoinedGroups `json:"JoinedGroups"`
}

func SearchData(w http.ResponseWriter, r *http.Request) {
	searchfromfrontend := r.URL.Query().Get("query")
    searchTerm := "%" + searchfromfrontend + "%"
	fmt.Println("", searchTerm)
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
	SELECT u.id, u.firstName, u.lastName, u.about, u.avatar, 
	       COUNT(f2.followerId) as followerCount
	FROM Users u
	LEFT JOIN Followers f2 ON f2.followedId = u.id
	WHERE u.id != ? AND u.id NOT IN (
		SELECT followedId FROM Followers WHERE followerId = ?
	)
	AND (u.firstName LIKE ? OR u.lastName LIKE ?)
	GROUP BY u.id
	`

	rows, err := dataB.SocialDB.Query(query, userID, userID, searchTerm, searchTerm)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var unfollowedUsers []NotfollowedUser
	for rows.Next() {
		var user NotfollowedUser
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.About, &user.Avatar, &user.FollowerCount)
		if err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		unfollowedUsers = append(unfollowedUsers, user)
	}

	query2 := `
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
	AND (g.title LIKE ? OR g.description LIKE ?)
	GROUP BY g.id, g.title, g.description
	`

	rows2, err := dataB.SocialDB.Query(query2, userID, searchTerm, searchTerm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows2.Close()
	var groups []NotJoinGroups
	for rows2.Next() {
		var g NotJoinGroups
		if err := rows2.Scan(&g.Id, &g.Title, &g.Description, &g.MembersCount, &g.PostCont); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}
		groups = append(groups, g)
	}


	query3 := `
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
	AND (g.title LIKE ? OR g.description LIKE ?)
	GROUP BY g.id, g.title, g.description
	`

	rows3, err := dataB.SocialDB.Query(query3, userID,searchTerm,searchTerm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows3.Close()
	var groups1 []JoinedGroups
	for rows3.Next() {
		var g JoinedGroups
		if err := rows3.Scan(&g.Id, &g.Title, &g.Description, &g.MembersCount, &g.PostCont); err != nil {
			fmt.Println("error to get groups", err)
			http.Error(w, "error to get groups", http.StatusInternalServerError)
			return
		}
		groups1 = append(groups1, g)
	}

	var databackend Data

	databackend.Notfollowed = unfollowedUsers
	databackend.UnJoinGroups = groups
	databackend.JoinedGroups = groups1
	fmt.Println("hi yousesf", databackend)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(databackend); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
