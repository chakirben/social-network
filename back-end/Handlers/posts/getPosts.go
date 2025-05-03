package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	dataB "socialN/dataBase"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	baseQuery := `
		SELECT
				p.id,
				p.title,
				p.image,
				p.content, 
				u.firstName,
				u.lastName, 
				(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1) AS likeCount,
				(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1) AS dislikeCount,
				(SELECT reactionType FROM PostReactions WHERE postId = p.id AND userId = ?) AS userReaction,
				p.createdAt
			FROM Posts p
			JOIN Users u ON p.creatorId = u.id
			WHERE 
				(
					(p.groupId IS NOT NULL AND EXISTS (
						SELECT 1 FROM GroupsMembers WHERE groupId = p.groupId AND memberId = ?
					))
					OR (p.privacy = 'public')
					OR (p.privacy = "almostPrivate" AND EXISTS (
						SELECT 1 FROM Followers WHERE followedId = p.creatorId AND followerId = ?
					))
					OR (p.privacy = "private" AND EXISTS (
						SELECT 1 FROM PostViewPermissions WHERE postId = p.id AND userId = ?
					))
				)
			ORDER BY p.createdAt DESC
			LIMIT 10 OFFSET ?;

	`
	rows, err := dataB.SocialDB.Query(baseQuery, 1, 1, 1, 1, offsetInt)
	if err != nil {
		log.Println("Error fetching posts:", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []map[string]interface{}{}
	for rows.Next() {
		var (
			id           int
			title        string
			image        sql.NullString
			content      string
			firstName    string
			lastName     string
			likeCount    sql.NullInt32
			dislikeCount sql.NullInt32
			userReaction sql.NullInt32
			createdAt    time.Time
		)
		err := rows.Scan(&id, &title, &image, &content, &firstName, &lastName, &likeCount, &dislikeCount, &userReaction, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		post := map[string]interface{}{
			"id":            id,
			"title":         title,
			"image":         image.String,
			"content":       content,
			"creator":       fmt.Sprintf("%s %s", firstName, lastName),
			"like_count":    likeCount.Int32,
			"dislike_count": dislikeCount.Int32,
			"user_reaction": userReaction.Int32,
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
