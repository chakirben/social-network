package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"socialN/app/auth"
	dataB "socialN/dataBase"
)

func GetCreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT
			p.id,
			p.content,
			p.image,
			u.avatar,
			u.firstName,
			u.lastName,
			(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1) AS likeCount,
			(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1) AS dislikeCount,
			(SELECT reactionType FROM postReactions WHERE postId = p.id AND userId = ?) AS userReaction,
			p.createdAt
		FROM Posts p
		JOIN Users u ON p.CreatorId = u.id
		WHERE p.creatorId = ?
		ORDER BY p.createdAt DESC
	`

	rows, err := dataB.SocialDB.Query(query, userID, userID)
	if err != nil {
		log.Println("Error fetching created posts:", err)
		http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []map[string]interface{}{}
	for rows.Next() {
		var postID int
		var content, firstName, lastName string
		var image *string
		var avatar *string
		var likeCount, dislikeCount int
		var userReaction sql.NullInt64 // ✅ Allow nullable
		var createdAt time.Time

		err := rows.Scan(&postID, &content, &image, &avatar, &firstName, &lastName,
			&likeCount, &dislikeCount, &userReaction, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Safely set userReaction value
		var reactionValue *int
		if userReaction.Valid {
			val := int(userReaction.Int64)
			reactionValue = &val
		}

		post := map[string]interface{}{
			"id":            postID,
			"content":       content,
			"image":         image,
			"creator":       fmt.Sprintf("%s %s", firstName, lastName),
			"avatar":        avatar,
			"like_count":    likeCount,
			"dislike_count": dislikeCount,
			"user_reaction": reactionValue,
			"created_at":    createdAt,
		}

		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
