package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Posts struct {
	Id           int       `json:"id"`
	Image        *string   `json:"image"`
	Content      string    `json:"content"`
	FirstName    string    `json:"creator"`
	LastName     string    `json:"groupid"`
	Avatar       string    `json:"avatar"`
	LikeCount    *int      `json:"like_count"`
	DislikeCount *int      `json:"dislike_count"`
	UserReaction *int      `json:"user_reaction"`
	CreatedAt    time.Time `json:"created_at"`
	GPTitle      string
}

func GetPostGroups(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session :(", http.StatusUnauthorized)
		return
	}

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
	var exists bool
	err = dataB.SocialDB.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", idInt).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	query := `SELECT 
	      P.id,
		  P.content,
		  P.image,
		  U.firstName,
          U.lastName,
          U.avatar,
          P.createdAt,
		  	(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1) AS likeCount,
			(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1) AS dislikeCount,
			(SELECT reactionType FROM postReactions WHERE postId = p.id AND userId = ?) AS userReaction
		FROM Posts P
		JOIN Users U ON P.creatorId = U.id
		WHERE groupId = ?
	`

	rows, err := dataB.SocialDB.Query(query, userID, idInt)
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
			&P.FirstName, &P.LastName,&P.Avatar, &P.CreatedAt, &P.LikeCount, &P.DislikeCount, &P.UserReaction); err != nil {
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
