package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"socialN/app/auth"
	database "socialN/dataBase"
)

type Message2 struct {
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id,omitempty"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Avatar     string    `json:"avatar"`
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	msgType := r.URL.Query().Get("type")
	switch msgType {
	case "private":
		handlePrivateMessages(w, r, userID)
	case "group":
		handleGroupMessages(w, r)
	default:
		http.Error(w, "Invalid or missing 'type' parameter. Use 'private' or 'group'", http.StatusBadRequest)
	}
}
func handlePrivateMessages(w http.ResponseWriter, r *http.Request, userID int) {
	otherIDStr := r.URL.Query().Get("other_id")
	if otherIDStr == "" {
		http.Error(w, "Missing 'other_id' parameter", http.StatusBadRequest)
		return
	}

	otherID, err := strconv.Atoi(otherIDStr)
	if err != nil {
		http.Error(w, "Invalid 'other_id' parameter", http.StatusBadRequest)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid 'offset' parameter", http.StatusBadRequest)
			return
		}
	}

	query := `
		SELECT 
			m.content, 
			m.sentAt,
			m.senderId,
			m.receiverId,
			u.firstName,
			u.lastName,
			u.avatar
		FROM messages m
		JOIN users u ON u.id = m.senderId
		WHERE 
			(m.senderId = ? AND m.receiverId = ?)
			OR (m.senderId = ? AND m.receiverId = ?)
		ORDER BY m.sentAt DESC
		LIMIT 12 OFFSET ?;
	`

	rows, err := database.SocialDB.Query(query, userID, otherID, otherID, userID, offset)
	if err != nil {
		fmt.Println("DB error:", err)
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message2
	for rows.Next() {
		var msg Message2
		if err := rows.Scan(&msg.Content, &msg.SentAt, &msg.SenderID, &msg.ReceiverID, &msg.FirstName, &msg.LastName, &msg.Avatar); err == nil {
			messages = append(messages, msg)
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].SentAt.Before(messages[j].SentAt)
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func handleGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "Missing 'group_id' parameter", http.StatusBadRequest)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	var err error
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid 'offset' parameter", http.StatusBadRequest)
			return
		}
	}

	query := `
		SELECT 
			m.content, 
			m.sentAt,
			m.senderId,
			0 as receiverId,
			u.firstName,
			u.lastName,
			u.avatar
		FROM GroupMessages m
		JOIN users u ON u.id = m.senderId
		WHERE m.groupId = ?
		ORDER BY m.sentAt DESC
		LIMIT 12 OFFSET ?;
	`

	rows, err := database.SocialDB.Query(query, groupID, offset)
	if err != nil {
		fmt.Println("DB error:", err)
		http.Error(w, "Error fetching group messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message2
	for rows.Next() {
		var msg Message2
		if err := rows.Scan(&msg.Content, &msg.SentAt, &msg.SenderID, &msg.ReceiverID, &msg.FirstName, &msg.LastName, &msg.Avatar); err == nil {
			messages = append(messages, msg)
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].SentAt.Before(messages[j].SentAt)
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}