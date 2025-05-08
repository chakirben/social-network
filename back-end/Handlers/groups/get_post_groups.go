package groups

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type IdOFgroup struct {
	Groupid int `json:"groupid"`
}

type Posts struct {
	Id           int
	Title        string
	Image        sql.NullString
	Content      string
	FirstName    string
	LastName     string
	LikeCount    sql.NullInt32
	DislikeCount sql.NullInt32
	UserReaction sql.NullInt32
	CreatedAt    time.Time
}

func GetPostGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}
	fmt.Println(userID)
	var Idgroup IdOFgroup
	err = json.NewDecoder(r.Body).Decode(&Idgroup)
	if err != nil {
		http.Error(w, "Invalid JSON :(", http.StatusBadRequest)
		return
	}
	if Idgroup.Groupid < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}
	query := `SELECT 
	      P.id,
		  P.title,
		  P.content,
		  P.image,
		  P.createdAt
		FROM Posts P
		JOIN Users U ON P.creatorId = U.id
		WHERE groupId = ?
	`

	rows, err := dataB.SocialDB.Query(query, Idgroup.Groupid)
	if err != nil {
		fmt.Println("--->",err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var allposts []Posts
	for rows.Next() {
		var P Posts
		if err := rows.Scan(&P.Id, &P.Title, &P.Content, &P.Image,
			 &P.CreatedAt); err != nil {
			fmt.Println("error to get posts groups", err)
			http.Error(w, "error to get posts groups", http.StatusInternalServerError)
			return
		}

		allposts = append(allposts, P)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allposts); err != nil {
		fmt.Println("JSON encode error", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
