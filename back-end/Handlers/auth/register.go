package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	db "socialN/dataBase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const maxUploadSize = 10 * 1024 * 1024

type User struct {
	ID          int
	Nickname    string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	DateOfBirth string
	Avatar      string
	About       string
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	nickname := r.FormValue("nickname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	dateOfBirth := r.FormValue("dob")
	about := r.FormValue("about")

	user := User{
		Nickname:    nickname,
		Email:       email,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		About:       about,
	}

	err = ValidateUser(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := checkUserExistence(user.Nickname, user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	avatarPath, err := handleAvatarUpload(r, w)
	if err != nil {
		http.Error(w, "Error handling avatar upload", http.StatusInternalServerError)
		fmt.Println("Error handling avatar upload:", err)
		return
	}
	user.Avatar = avatarPath

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := insertUser(db.SocialDB, user, string(hashedPassword)); err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	err = db.SocialDB.QueryRow("SELECT id FROM Users WHERE email = ?", user.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		return
	}

	sessionID, expiration, err := CreateSession(db.SocialDB, user.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func handleAvatarUpload(r *http.Request, w http.ResponseWriter) (string, error) {
	file, header, err := r.FormFile("avatar")
	if err == http.ErrMissingFile {
		return "", nil
	} else if err != nil {
		http.Error(w, "Error handling avatar upload", http.StatusInternalServerError)
		return "", err
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		http.Error(w, "Avatar size too large", http.StatusBadRequest)
		return "", fmt.Errorf("avatar size too large")
	}

	allowedImageTypes := map[string]bool{"image/jpeg": true, "image/png": true, "image/gif": true}
	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		http.Error(w, "Invalid avatar file type", http.StatusBadRequest)
		return "", fmt.Errorf("invalid avatar file type")
	}

	uuidStr := uuid.New().String()
	filename := fmt.Sprintf("%s_%s", uuidStr, header.Filename)
	localPath := filepath.Join("uploads", filename)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0o755)
	}

	dst, err := os.Create(localPath)
	if err != nil {
		http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
		return "", err
	}

	publicPath := fmt.Sprintf("/uploads/%s", filename)
	return publicPath, nil
}

func checkUserExistence(nickname, email string) error {
	fmt.Println(nickname, email)
	if nickname != "" {
		var existingNickname string
		err := db.SocialDB.QueryRow("SELECT nickname FROM Users WHERE nickname = ?", nickname).Scan(&existingNickname)
		if err == nil {
			return fmt.Errorf("nickname already exists")
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("error checking nickname existence: %w", err)
		}
	}

	var existingEmail string
	err := db.SocialDB.QueryRow("SELECT email FROM Users WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		return fmt.Errorf("email already exists")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking email existence: %w", err)
	}
	return nil
}

func insertUser(db *sql.DB, user User, hashedPassword string) error {
	query := `
		INSERT INTO Users (nickname, email, password, firstName, lastName, dateOfBirth, avatar, about, accountType)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'public')
	`
	_, err := db.Exec(query,
		user.Nickname,
		user.Email,
		hashedPassword,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Avatar,
		user.About,
	)
	if err != nil {
		log.Println("Error inserting user into database:", err)
	}
	return err
}

func ValidateUser(user User) error {
	fmt.Println(user)
	if strings.TrimSpace(user.Email) == "" ||
		strings.TrimSpace(user.Password) == "" ||
		strings.TrimSpace(user.FirstName) == "" ||
		strings.TrimSpace(user.LastName) == "" ||
		strings.TrimSpace(user.DateOfBirth) == "" {
		return errors.New("please fill in all mandatory fields")
	}

	dob, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		return errors.New("invalid date of birth format (YYYY-MM-DD)")
	}
	age := time.Now().Year() - dob.Year()
	if age < 16 {
		return errors.New("you must be at least 16 years old")
	}
	if user.Nickname != "" {
		if len(user.Nickname) >= 24 {
			return errors.New("nickname should be less than or equal to 24 characters")
		}
		nicknameRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
		if !nicknameRegex.MatchString(user.Nickname) {
			return errors.New("nickname can only contain letters, numbers, underscores, and periods")
		}
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !nameRegex.MatchString(user.FirstName) {
		return errors.New("first name must contain only letters")
	}
	if !nameRegex.MatchString(user.LastName) {
		return errors.New("last name must contain only letters")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}
