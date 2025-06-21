package posts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"socialN/app/auth"
	cmnts "socialN/app/comments"
	dataB "socialN/dataBase"
)

type PostResponse struct {
	Id           int       `json:"id"`
	Image        *string   `json:"image"`
	Content      string    `json:"content"`
	Creator      string    `json:"creator"`
	Avatar       string    `json:"avatar"`
	GroupId      string    `json:"groupid"`
	LikeCount    *int      `json:"like_count"`
	DislikeCount *int      `json:"dislike_count"`
	UserReaction *int      `json:"user_reaction"`
	CreatedAt    time.Time `json:"created_at"`
	GPTitle      string
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	privacy := r.FormValue("privacy")
	if privacy != "public" && privacy != "almostPrivate" && privacy != "private" && privacy != "inGroup" {
		http.Error(w, "Invalid post privacy", http.StatusBadRequest)
		return
	}

	selectedUsers := r.MultipartForm.Value["selectedUsers"]
	if privacy == "private" && len(selectedUsers) == 0 {
		http.Error(w, "No selected audience for private post", http.StatusBadRequest)
		return
	}

	var imagePath string
	file, _, err := r.FormFile("image")
	if err == nil && file != nil {
		imagePath, err = cmnts.SaveAvatar(file)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	}

	groupIDStr := r.FormValue("groupId")
	var groupInt int
	var result any
	fmt.Println(groupIDStr, privacy)
	if privacy == "inGroup" && groupIDStr != "" {
		groupInt, err = strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
		result, err = dataB.SocialDB.Exec(`
			INSERT INTO Posts (content, image, privacy, groupId, creatorId)
			VALUES (?, ?, ?, ?, ?)`,
			content, imagePath, privacy, groupInt, userID,
		)
	} else {
		result, err = dataB.SocialDB.Exec(`
			INSERT INTO Posts (content, image, privacy, creatorId)
			VALUES (?, ?, ?, ?)`,
			content, imagePath, privacy, userID,
		)
	}
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		fmt.Println("Insert error:", err)
		return
	}

	postID, err := result.(interface {
		LastInsertId() (int64, error)
	}).LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}

	if privacy == "private" {
		for _, idStr := range selectedUsers {
			_, err := dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO PostViewPermissions (postId, userId)
				VALUES (?, ?)`, postID, idStr)
			if err != nil {
				http.Error(w, "Permission insert error", http.StatusInternalServerError)
			}
		}
	}

	baseQuery := `
		SELECT
			p.image,
			p.content, 
			u.firstName || ' ' || u.lastName,
			COALESCE(u.avatar, ''),
			(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1) AS likeCount,
			(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1) AS dislikeCount,
			(SELECT reactionType FROM PostReactions WHERE postId = p.id AND userId = ?) AS userReaction,
			p.createdAt
		FROM Posts p
		JOIN Users u ON p.creatorId = u.id
		WHERE 
			(
				(p.creatorId = ?) OR
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
		AND p.id = ?
	`

	var post PostResponse

	err = dataB.SocialDB.QueryRow(baseQuery, userID, userID, userID, userID, userID, postID).Scan(
		&post.Image,
		&post.Content,
		&post.Creator,
		&post.Avatar,
		&post.LikeCount,
		&post.DislikeCount,
		&post.UserReaction,
		&post.CreatedAt,
	)
	if err != nil {
		log.Println("Error fetching posts:", err)
		fmt.Println(err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	post.Id = int(postID)
	post.GroupId = groupIDStr

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
