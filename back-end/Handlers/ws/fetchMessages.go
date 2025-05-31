package ws

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"socialN/Handlers/auth"
	database "socialN/dataBase"
)

type Message2 struct {
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id,omitempty"`
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	msgType := r.URL.Query().Get("type")

	var messages []Message2

	switch msgType {
	case "private":
		otherID := r.URL.Query().Get("other_id")
		if otherID == "" {
			http.Error(w, "Missing 'other_id' parameter", http.StatusBadRequest)
			return
		}

		query := `
			SELECT 
				m.content, 
				m.sent_at,
				m.sender_id,
				m.receiver_id
			FROM messages m
			WHERE 
				(m.sender_id = ? AND m.receiver_id = ?)
				OR (m.sender_id = ? AND m.receiver_id = ?)
			ORDER BY m.sent_at DESC
			LIMIT 50;
		`

		rows, err := database.SocialDB.Query(query, userID, otherID, otherID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var msg Message2
			err := rows.Scan(&msg.Content, &msg.SentAt, &msg.SenderID, &msg.ReceiverID)
			if err != nil {
				continue
			}
			messages = append(messages, msg)
		}

	case "group":
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "Missing 'group_id' parameter", http.StatusBadRequest)
			return
		}

		query := `
			SELECT 
				gm.content,
				gm.sent_at,
				gm.senderId
			FROM GroupMessages gm
			WHERE gm.groupId = ?
			ORDER BY gm.sent_at DESC
			LIMIT 50;
		`

		rows, err := database.SocialDB.Query(query, groupID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var msg Message2
			err := rows.Scan(&msg.Content, &msg.SentAt, &msg.SenderID)
			if err != nil {
				continue
			}
			messages = append(messages, msg)
		}

	default:
		http.Error(w, "Invalid or missing 'type' parameter. Use 'private' or 'group'", http.StatusBadRequest)
		return
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].SentAt.Before(messages[j].SentAt)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
