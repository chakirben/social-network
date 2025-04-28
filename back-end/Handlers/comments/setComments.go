package comment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	dataB "socialN/dataBase"
)

type CommentResponse struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	Content    string    `json:"content"`
	Image      string    `json:"image"`
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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	postID := r.FormValue("postId")
	image, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving image", http.StatusInternalServerError)
		return
	}

	var imagePath string
	if image != nil {

		imagePath, err = saveAvatar(image)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	}

	// userID, err := auth.ValidateSession(r, dataB.SocialDB)
	// if err != nil {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	fmt.Println(postID , content ,  imagePath)
	result, err := dataB.SocialDB.Exec(`
		INSERT INTO Comments (postId, userId, content, image)
		VALUES (?, ?, ?, ?)`,
		postID, 1, content, imagePath,
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
		SELECT Comments.id, Users.nickName, Users.firstName, Users.lastName, Comments.createdAt, Comments.content, Comments.image
		FROM Comments
		JOIN Users ON Comments.userId = Users.id
		WHERE Comments.id = ?`, lastID)

	var commentResponse CommentResponse
	err = row.Scan(&commentResponse.Id, &commentResponse.Username, &commentResponse.First_name, &commentResponse.Last_name, &commentResponse.CreatedAt, &commentResponse.Content, &commentResponse.Image)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve created comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commentResponse)
}

func saveAvatar(file io.Reader) (string, error) {
	// Save the uploaded image file to the server
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_avatar.jpg", timestamp)
	avatarPath := filepath.Join("uploads", filename)

	// Create upload folder if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0o755)
	}

	// Save the file
	dst, err := os.Create(avatarPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return avatarPath, nil
}
