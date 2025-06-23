package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"socialN/app/auth"
	database "socialN/dataBase"
)

type Discussion struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	IsGroup            bool       `json:"isGroup"`
	Avatar             *string    `json:"avatar"` // Avatar URL or path, nullable
	LastMessageContent *string    `json:"lastMessageContent"`
	LastMessageSentAt  *time.Time `json:"lastMessageSentAt"`
}

func GetAllDiscussionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in Disc")

	userID, err := auth.ValidateSession(r, database.SocialDB)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var discussions []Discussion

	// Private conversations query
	privateQuery := `
	SELECT 
		u.id,
		u.firstName,
		u.lastName,
		u.avatar,
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
	WHERE u.id != ?
	AND EXISTS (
 		SELECT 1 FROM Messages m
  		WHERE 
    		(m.senderId = ? AND m.receiverId = u.id)
    		OR (m.senderId = u.id AND m.receiverId = ?)
	)
	`

	privateRows, err := database.SocialDB.Query(privateQuery,
		userID, userID, // for last_message_content subquery
		userID, userID, // for last_message_sent_at subquery
		userID,         // for WHERE u.id != ?
		userID, userID, // for EXISTS subquery
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Private query error:", err)
		return
	}
	defer privateRows.Close()

	for privateRows.Next() {
		var id int
		var firstName, lastName string
		var avatar sql.NullString
		var content sql.NullString
		var sentAt sql.NullTime

		if err := privateRows.Scan(&id, &firstName, &lastName, &avatar, &content, &sentAt); err != nil {
			fmt.Println("Private row scan error:", err)
			continue
		}

		fullName := firstName + " " + lastName

		var avatarPtr *string
		if avatar.Valid {
			avatarPtr = &avatar.String
		}

		var lastMsgContent *string
		var lastMsgSentAt *time.Time
		if content.Valid {
			lastMsgContent = &content.String
		}
		if sentAt.Valid {
			lastMsgSentAt = &sentAt.Time
		}

		discussions = append(discussions, Discussion{
			ID:                 id,
			Name:               fullName,
			IsGroup:            false,
			Avatar:             avatarPtr,
			LastMessageContent: lastMsgContent,
			LastMessageSentAt:  lastMsgSentAt,
		})
	}

	// Group conversations query (no avatar for groups, only title)
	groupQuery := `
	SELECT 
		g.id,
		g.title,
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
	WHERE gm.memberId = ?
	  AND EXISTS (
		  SELECT 1 FROM GroupMessages m
		  WHERE m.groupId = g.id
	  );
	`

	groupRows, err := database.SocialDB.Query(groupQuery, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Group query error:", err)
		return
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var groupId int
		var groupTitle string
		var content sql.NullString
		var sentAt sql.NullTime

		if err := groupRows.Scan(&groupId, &groupTitle, &content, &sentAt); err != nil {
			fmt.Println("Group row scan error:", err)
			continue
		}

		var lastMsgContent *string
		var lastMsgSentAt *time.Time
		if content.Valid {
			lastMsgContent = &content.String
		}
		if sentAt.Valid {
			lastMsgSentAt = &sentAt.Time
		}

		discussions = append(discussions, Discussion{
			ID:                 groupId,
			Name:               groupTitle,
			IsGroup:            true,
			Avatar:             nil, // no avatar for groups
			LastMessageContent: lastMsgContent,
			LastMessageSentAt:  lastMsgSentAt,
		})
	}

	// Sort discussions by last message time descending (nil times go last)
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
