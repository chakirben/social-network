package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	database "socialN/dataBase"
)

type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
type MessagesRequest struct {
	ReceiverID int `json:"receiver_id"`
	// limit and offset
	Options string `json:"options"`
}

// Insert message
func InsertMessage(message Message) error {
	//  check if the message is valid
	if message.SenderID == 0 || message.ReceiverID == 0 || message.Content == "" {
		return errors.New("missing required fields")
	}
	_, err := database.SocialDB.Exec(`INSERT INTO Messages (senderId, receiverId, content, sentAt) VALUES (?, ?, ?, ?)`,
		message.SenderID, message.ReceiverID, message.Content, message.CreatedAt)
	return err
}

// Get messages and users
func GetMessagesAndUsers() ([]Message, error) {
	rows, err := database.SocialDB.Query(`SELECT m.id, m.senderId, m.receiverId, m.content, m.sentAt, u.firstName, u.lastName, u.avatar FROM Messages m JOIN Users u ON m.senderId = u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	i := 1
	for rows.Next() && i < 10 {
		var message Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.Content, &message.Content, &message.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
		i++
	}
	return messages, nil
}

// Fetch  More  messages
func FetchMoreMessages (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var options MessagesRequest
	json.NewDecoder(r.Body).Decode(&options)
	opt := strings.Split(options.Options, ":")
	if len(opt) != 2 {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}
	offset, _ := strconv.Atoi(opt[0])
	limit, _ := strconv.Atoi(opt[1])
	rows, err := database.SocialDB.Query(`SELECT m.id, m.senderId, m.receiverId, m.content, m.sentAt, u.firstName, u.lastName, u.avatar FROM Messages m JOIN Users u ON m.senderId = u.id LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.Content, &message.Content, &message.Content)
		if err != nil {
			http.Error(w, "Error fetching messages", http.StatusInternalServerError)
			return
		}
		if message.ReceiverID != options.ReceiverID {
			continue
		}
		messages = append(messages, message)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
