package comment

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

type Comment struct {
	ID           int     `json:"id"`
	Content      string  `json:"content"`
	Image        *string `json:"image"`
	FirstName    string  `json:"firstName"` 
	LastName     string  `json:"lastName"`
	Avatar       *string `json:"avatar"`
	CreatedAt    string  `json:"createdAt"`
	LikeCount    int     `json:"likeCount"`
	DislikeCount int     `json:"dislikeCount"`
	UserReaction *int    `json:"userReaction"`
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		http.Error(w, "Missing post_id query parameter", http.StatusBadRequest)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	commentRows, err := dataB.SocialDB.Query(`
		SELECT 
			c.id,
			c.content, 
			c.image,
			u.firstName, 
			u.lastName, 
			u.avatar,
			c.createdAt,
			(SELECT COUNT(*) FROM commentReactions WHERE commentId = c.id AND reactionType = 1),
			(SELECT COUNT(*) FROM commentReactions WHERE commentId = c.id AND reactionType = -1),
			(SELECT reactionType FROM commentReactions WHERE commentId = c.id AND userId = ?)
		FROM Comments c
		JOIN Users u ON c.userId = u.id
		WHERE c.postId = ?
		ORDER BY c.createdAt DESC
	`,userID, postID)
	if err != nil {
		log.Println("Error fetching comments:", err)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer commentRows.Close()

	var comments []Comment

	for commentRows.Next() {
		var comment Comment
		err := commentRows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.Image,
			&comment.FirstName,
			&comment.LastName,
			&comment.Avatar,
			&comment.CreatedAt,
			&comment.LikeCount,
			&comment.DislikeCount,
			&comment.UserReaction,
		)
		if err != nil {
			log.Println("Error scanning comment:", err)
			continue
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "Failed to encode comments to JSON", http.StatusInternalServerError)
	}
}
