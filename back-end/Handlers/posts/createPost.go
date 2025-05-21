package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"socialN/Handlers/auth"
	cmnts "socialN/Handlers/comments"
	dataB "socialN/dataBase"
)

type PostResponse struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Image   string `json:"image"`
	Privacy string `json:"privacy"`
	GroupID int   `json:"groupId,omitempty"`
	Creator int    `json:"creatorId"`
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
	if privacy != "public" && privacy != "almostPrivate" && privacy != "private" &&  privacy != "inGroup" {
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

	// Handle groupId if present
	groupIDStr := r.FormValue("groupId")
	var groupInt int
	var result any
	fmt.Println(groupIDStr , privacy )
	if privacy =="inGroup" && groupIDStr != "" {

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
			content,imagePath, privacy, userID,
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

	// Insert view permissions for private posts
	if privacy == "private" {
		for _, idStr := range selectedUsers {
			_, err := dataB.SocialDB.Exec(`
				INSERT OR IGNORE INTO PostViewPermissions (postId, userId)
				VALUES (?, ?)`, postID, idStr)
			if err != nil {
				http.Error(w,"Permission insert error:" , http.StatusInternalServerError)
			}
		}
	}

	response := PostResponse{
		ID:      postID,
		Content: content,
		Image:   imagePath,
		Privacy: privacy,
		Creator: userID,
		GroupID: groupInt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
