package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dataB "socialN/dataBase"
)

func GetCreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Invalid session", http.StatusUnauthorized)
	// 	return
	// }

	query := `
		SELECT
			p.id,
			p.title,
			p.content,
			p.image,
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
	rows, err := dataB.SocialDB.Query(query, 2, 2)
	if err != nil {
		log.Println("Error fetching created posts:", err)
		http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []map[string]interface{}{}

	for rows.Next() {
		var postID int
		var title, content, firstName, lastName string
		var image sql.NullString
		var likeCount, dislikeCount, userReaction sql.NullInt32
		var createdAt time.Time

		err := rows.Scan(&postID, &title, &content,&image, &firstName, &lastName, &likeCount, &dislikeCount, &userReaction, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		post := map[string]interface{}{
			"id":           postID,
			"title":        title,
			"content":      content,
			"image":      image,
			"creator":      fmt.Sprintf("%s %s", firstName, lastName),
			"likeCount":    likeCount.Int32,
			"dislikeCount": dislikeCount.Int32,
			"userReaction": userReaction.Int32,
			"createdAt":    createdAt,
		}

		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
