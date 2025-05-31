package ws

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"socialN/Handlers/auth"
	database "socialN/dataBase"
)

type Discussion struct {
	ID                 int        `json:"id"`   
	Name               string     `json:"name"`
	IsGroup            bool       `json:"isGroup"`
	LastMessageContent *string    `json:"lastMessageContent"`
	LastMessageSentAt  *time.Time `json:"lastMessageSentAt"`
}

func GetAllDiscussionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var discussions []Discussion
	privateQuery := `
	SELECT 
		u.id,
		u.nickname,
		(
			SELECT m.content
			FROM Messages m
			WHERE 
				(m.senderId = ? AND m.receiverId = u.id)
				OR (m.senderId = u.id AND m.receiverId = ?)
			ORDER BY m.sentAt DESC
			LIMIT 1
		) AS last_message_content,
		(
			SELECT m.sentAt
			FROM Messages m
			WHERE 
				(m.senderId = ? AND m.receiverId = u.id)
				OR (m.senderId = u.id AND m.receiverId = ?)
			ORDER BY m.sentAt DESC
			LIMIT 1
		) AS last_message_sent_at
	FROM users u
	WHERE u.id != ?;
	`

	privateRows, err := database.SocialDB.Query(privateQuery, userID, userID, userID, userID, userID)
	if err == nil {
		defer privateRows.Close()
		for privateRows.Next() {
			var id int
			var nickname string
			var content sql.NullString
			var sentAt sql.NullTime

			if err := privateRows.Scan(&id, &nickname, &content, &sentAt); err == nil {
				if content.Valid && sentAt.Valid {
					discussions = append(discussions, Discussion{
						ID:                 id,
						Name:               nickname,
						IsGroup:            false,
						LastMessageContent: &content.String,
						LastMessageSentAt:  &sentAt.Time,
					})
				}
			}
		}
	}

	groupQuery := `
	SELECT 
		g.id,
		g.name,
		(
			SELECT gm.content
			FROM GroupMessages gm
			WHERE gm.groupId = g.id
			ORDER BY gm.sentAt DESC
			LIMIT 1
		) AS last_message_content,
		(
			SELECT gm.sentAt
			FROM GroupMessages gm
			WHERE gm.groupId = g.id
			ORDER BY gm.sentAt DESC
			LIMIT 1
		) AS last_message_sent_at
	FROM groups g
	JOIN GroupsMembers gm ON gm.groupId = g.id
	WHERE gm.memberId = ?;
	`

	groupRows, err := database.SocialDB.Query(groupQuery, userID)
	if err == nil {
		defer groupRows.Close()
		for groupRows.Next() {
			var groupId int
			var groupName string
			var content sql.NullString
			var sentAt sql.NullTime

			if err := groupRows.Scan(&groupId, &groupName, &content, &sentAt); err == nil {
				if content.Valid && sentAt.Valid {
					discussions = append(discussions, Discussion{
						ID:                 groupId,
						Name:               groupName,
						IsGroup:            true,
						LastMessageContent: &content.String,
						LastMessageSentAt:  &sentAt.Time,
					})
				}
			}
		}
	}

	sort.Slice(discussions, func(i, j int) bool {
		if discussions[i].LastMessageSentAt == nil {
			return false
		}
		if discussions[j].LastMessageSentAt == nil {
			return true
		}
		return discussions[i].LastMessageSentAt.After(*discussions[j].LastMessageSentAt)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discussions)
}
