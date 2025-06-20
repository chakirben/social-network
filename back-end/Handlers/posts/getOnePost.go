package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var id int
	var content, firstName, lastName, avatar string
	var image sql.NullString
	var likeCount, dislikeCount, userReaction sql.NullInt32
	var createdAt time.Time

	row := dataB.SocialDB.QueryRow(`
        SELECT
            p.id,
            p.image,
            p.content,
            u.firstName,
            u.lastName,
            u.avatar,
            (SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1),
            (SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1),
            (SELECT reactionType FROM postReactions WHERE postId = p.id AND userId = ?),
            p.createdAt
        FROM Posts p
        JOIN Users u ON p.creatorId = u.id
        WHERE
            p.id = ?
            AND (
                (p.creatorId = ?)
                OR (p.groupId IS NOT NULL AND EXISTS (
                    SELECT 1 FROM GroupsMembers WHERE groupId = p.groupId AND memberId = ?
                ))
                OR (p.privacy = 'public')
                OR (p.privacy = 'almostPrivate' AND EXISTS (
                    SELECT 1 FROM Followers WHERE followedId = p.creatorId AND followerId = ?
                ))
                OR (p.privacy = 'private' AND EXISTS (
                    SELECT 1 FROM PostViewPermissions WHERE postId = p.id AND userId = ?
                ))
            )
    `, userID, postID, userID, userID, userID, userID, userID)

	err = row.Scan(&id, &image, &content, &firstName, &lastName, &avatar, &likeCount, &dislikeCount, &userReaction, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized or post not found", http.StatusUnauthorized)
			return
		}
		log.Println("Error fetching post:", err)
		http.Error(w, "Error fetching post", http.StatusInternalServerError)
		return
	}

	post := map[string]interface{}{
		"id":            id,
		"content":       content,
		"image":         image.String,
		"creator":       fmt.Sprintf("%s %s", firstName, lastName),
		"avatar":        avatar,
		"like_count":    likeCount.Int32,
		"dislike_count": dislikeCount.Int32,
		"user_reaction": userReaction.Int32,
		"created_at":    createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}