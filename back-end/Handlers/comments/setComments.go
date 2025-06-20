package comment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"

	"github.com/google/uuid"
)

type CommentResponse struct {
	Id           int       `json:"id"`
	Avatar       string    `json:"avatar"`
	First_name   string    `json:"firstName"`
	Last_name    string    `json:"lastName"`
	CreatedAt    time.Time `json:"createdAt"`
	Content      string    `json:"content"`
	Image        string    `json:"image"`
	MineReaction int       `json:"userReaction"`
}

type CommentRequest struct {
	PostID  int    `json:"postId"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

func SetCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Error retrieving content of the comment", http.StatusBadRequest)
		return
	}
	postID := r.FormValue("postId")
	image, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving image", http.StatusInternalServerError)
		return
	}

	var imagePath string
	if image != nil {
		imagePath, err = SaveAvatar(image)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	}

	result, err := dataB.SocialDB.Exec(`
		INSERT INTO Comments (postId, userId, content, image)
		VALUES (?, ?, ?, ?)`,
		postID, userID, content, imagePath,
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get last inserted ID", http.StatusInternalServerError)
		return
	}

	row := dataB.SocialDB.QueryRow(`
		SELECT Comments.id, Users.avatar, Users.firstName, Users.lastName, Comments.createdAt, Comments.content, Comments.image
		FROM Comments
		JOIN Users ON Comments.userId = Users.id
		WHERE Comments.id = ?`, lastID)

	var commentResponse CommentResponse
	commentResponse.MineReaction = 0

	err = row.Scan(&commentResponse.Id, &commentResponse.Avatar, &commentResponse.First_name, &commentResponse.Last_name, &commentResponse.CreatedAt, &commentResponse.Content, &commentResponse.Image)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve created comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commentResponse)
}

func SaveAvatar(file io.Reader) (string, error) {
	uuidStr := uuid.New().String()
	filename := fmt.Sprintf("%s_avatar.jpg", uuidStr)

	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, 0o755); err != nil {
			return "", err
		}
	}

	localPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("/uploads/%s", filename)
	return publicURL, nil
}
