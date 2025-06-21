package groups

import (
	"encoding/json"
	"log"
	"net/http"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

type Grouppp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GroupResponse struct {
	Id           int
	Title        string
	Description  string
	MembersCount int
	PostCont     int
}

func Creat_Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

	var group Grouppp
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}

	if group.Title == "" || group.Description == "" {
		http.Error(w, "Title and Description are required.", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO groups (title, description, adminId) VALUES (?, ?, ?)`
	res, err := dataB.SocialDB.Exec(query, group.Title, group.Description, userID)
	if err != nil {
		log.Println("Error inserting group:", err)
		http.Error(w, "Failed to create group.", http.StatusInternalServerError)
		return
	}

	lastIDgroup, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get last inserted ID.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	queryMember := `INSERT INTO GroupsMembers (memberId, groupId) VALUES (?, ?)`
	_, err = dataB.SocialDB.Exec(queryMember, userID, lastIDgroup)
	if err != nil {
		log.Println("Error inserting group member:", err)
		http.Error(w, "Failed to create group member.", http.StatusInternalServerError)
		return
	}

	querySelect := `
	SELECT 
		g.id, 
		g.title, 
		g.description, 
		COUNT(DISTINCT gm2.memberId) AS members_count,
		COUNT(DISTINCT P.groupId) AS post_count
	FROM Groups g
	LEFT JOIN GroupsMembers gm2 ON g.id = gm2.groupId
	LEFT JOIN Posts P ON P.groupId = g.id
	WHERE g.id = ?
	GROUP BY g.id, g.title, g.description
	`

	row := dataB.SocialDB.QueryRow(querySelect, lastIDgroup)

	var groupResp GroupResponse
	err = row.Scan(&groupResp.Id, &groupResp.Title, &groupResp.Description, &groupResp.MembersCount, &groupResp.PostCont)
	if err != nil {
		log.Println("Error fetching created group:", err)
		http.Error(w, "Failed to fetch created group.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(groupResp); err != nil {
		log.Println("JSON encode error:", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
