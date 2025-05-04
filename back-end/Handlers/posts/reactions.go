package posts

import (
	"fmt"
	"net/http"
	"strconv"

	database "socialN/dataBase"
)

type ReactionType struct {
	Id     int `json:"id"`
	Postid int `json:"postid"`
	UserId int `json:"userid"`
	Type   int `json:"type"`
}

func Reactpost(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	idStr := r.URL.Query().Get("id")
	postidStr := r.URL.Query().Get("postid")
	useridStr := r.URL.Query().Get("userid")
	typeStr := r.URL.Query().Get("type")

	// Convert to integers
	id, err := strconv.Atoi(idStr)
	postid, err2 := strconv.Atoi(postidStr)
	userid, err3 := strconv.Atoi(useridStr)
	reactionType, err4 := strconv.Atoi(typeStr)

	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// Validate reaction type
	if reactionType != -1 && reactionType != 0 && reactionType != 1 {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	// Insert into the database
	stmt, err := database.SocialDB.Prepare(`
	INSERT OR REPLACE INTO PostReactions (id, postId, userId, reactionType) VALUES (?, ?, ?, ?)
`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, postid, userid, reactionType)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Reaction added successfully")
}
