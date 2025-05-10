package groups

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	dataB "socialN/dataBase"
)

type Posts struct {
	Id           int            `json:"id"`
	Image        sql.NullString `json:"image"`
	Content      string         `json:"content"`
	FirstName    string         `json:"creator"`
	LastName     string         `json:"groupid"`
	LikeCount    sql.NullInt32  `json:"likeCount"`
	DislikeCount sql.NullInt32  `json:"dislikeCount"`
	UserReaction sql.NullInt32  `json:"userReaction"`
	CreatedAt    time.Time      `json:"createdAt"`
}

func GetPostGroups(w http.ResponseWriter, r *http.Request) {
	// Get the `id` parameter from the URL query
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "Missing group ID", http.StatusBadRequest)
		return
	}

	// You can now convert it to int if needed
	idInt, err := strconv.Atoi(groupID)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	
	if idInt < 1 {
		http.Error(w, "Invalid request :(", http.StatusBadRequest)
		return
	}
	query := `SELECT 
	      P.id,
		  P.content,
		  P.image,
		  U.firstName,
          U.lastName,
          P.createdAt
		FROM Posts P
		JOIN Users U ON P.creatorId = U.id
		WHERE groupId = ?
	`

	rows, err := dataB.SocialDB.Query(query, idInt)
	if err != nil {
		fmt.Println("--->", err)
		http.Error(w, "error to get my groups :(", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var allposts []Posts
	for rows.Next() {
		var P Posts
		if err := rows.Scan(&P.Id, &P.Content, &P.Image,
			&P.FirstName, &P.LastName, &P.CreatedAt); err != nil {
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
