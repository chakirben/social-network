package followers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"socialN/app/auth"
	"socialN/app/ws"
	dataB "socialN/dataBase"

	"github.com/gorilla/websocket"
)

type FollowInfo struct {
	Follower_session string `json:"follower_session"`
	Followed_id      string `json:"followed_id"`
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FollowHandler called")
	userID, err := auth.ValidateSession(r, dataB.SocialDB)
	if err != nil {
		fmt.Println("Unauthorized access:", err, userID)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetIDStr := r.URL.Query().Get("id")
	action := r.URL.Query().Get("action")

	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil || (action != "follow" && action != "unfollow" && action != "cancel_request") {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if targetID == userID {
		http.Error(w, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	db := dataB.SocialDB
	var accountType string
	err = db.QueryRow("SELECT accountType FROM Users WHERE id = ?", targetID).Scan(&accountType)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	switch action {
	case "follow":
		if accountType == "private" {
			var exists bool
			err = db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM Notifications
            WHERE senderId = ? AND receiverId = ? AND type = 'follow_request' AND status = 'pending'
        )
    `, userID, targetID).Scan(&exists)
			if err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}

			if !exists {
				// Insert notification
				res, err := db.Exec(`
            INSERT INTO Notifications (senderId, receiverId, type, status)
            VALUES (?, ?, 'follow_request', 'pending')
        `, userID, targetID)
				if err != nil {
					http.Error(w, "Couldn't send follow request", http.StatusInternalServerError)
					return
				}

				notificationId, _ := res.LastInsertId()

				// Fetch sender info for the WebSocket message
				var firstName, lastName, avatar string
				err = db.QueryRow("SELECT firstName, lastName, COALESCE(avatar, '') FROM Users WHERE id = ?", userID).
					Scan(&firstName, &lastName, &avatar)
				if err == nil {
					// Send WS notification to target user
					ws.ConnMu.Lock()
					conns := ws.Connections[targetID]
					ws.ConnMu.Unlock()

					for _, conn := range conns {
						if conn != nil {
							SendMsg(conn, NotificationMessage{
								ID:               notificationId,
								Type:             "Notification",
								NotificationType: "follow",
								FirstName:        firstName,
								LastName:         lastName,
								Avatar:           avatar,
							})
						}
					}
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("cancel_request"))
			return
		} else {
			_, err := db.Exec(`
                INSERT OR IGNORE INTO Followers (followerId, followedId)
                VALUES (?, ?)
            `, userID, targetID)
			if err != nil {
				http.Error(w, "Couldn't follow user", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("unfollow"))
			return
		}

	case "cancel_request":
		_, err := db.Exec(`
            DELETE FROM Notifications
            WHERE senderId = ? AND receiverId = ? AND type = 'follow_request' AND status = 'pending'
        `, userID, targetID)
		if err != nil {
			http.Error(w, "Couldn't cancel follow request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("follow"))
		return

	case "unfollow":
		_, err := db.Exec(`
            DELETE FROM Followers
            WHERE followerId = ? AND followedId = ?
        `, userID, targetID)
		if err != nil {
			http.Error(w, "Couldn't unfollow user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("follow"))
		return
	}

	// fallback
	http.Error(w, "Unknown error", http.StatusInternalServerError)
}

type NotificationMessage struct {
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	NotificationType string `json:"notificationType"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Avatar           string `json:"avatar"`
}

func SendMsg(conn *websocket.Conn, msg NotificationMessage) {
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshaling notification:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

type acceptFollow struct {
	FollowerID      int    `json:"follower_id"`
	FollowedSession string `json:"followed_session"`
}

func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	var followInfo acceptFollow

	err := json.NewDecoder(r.Body).Decode(&followInfo)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	var followedid int
	followerid := followInfo.FollowerID

	errq2 := dataB.SocialDB.QueryRow(`SELECT userId FROM Sessions WHERE id=?`, followInfo.FollowedSession).Scan(&followedid)
	if errq2 != nil {
		fmt.Println("Error get ID:", errq2)
		return
	}

	_, exec_err := dataB.SocialDB.Exec(`INSERT INTO Followers (followerId, followedId) VALUES (?,?)`, followerid, followedid)
	if exec_err != nil {
		fmt.Println("Error Insert into db:", exec_err)
		return
	}

	stmt, err := dataB.SocialDB.Prepare("UPDATE Notifications SET status = ? WHERE senderId = ? AND receiverId = ?")
	if err != nil {
		fmt.Println("Error change elem in db:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("accepted", followerid, followedid)
	if err != nil {
		fmt.Println("Error change elem in db:", err)
		return
	}

	response := map[string]interface{}{
		"status": "follow accepted successfuly",
	}
	json.NewEncoder(w).Encode(response)
}

func DeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	var followInfo acceptFollow

	err := json.NewDecoder(r.Body).Decode(&followInfo)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	var followedid int
	followerid := followInfo.FollowerID

	errq2 := dataB.SocialDB.QueryRow(`SELECT userId FROM Sessions WHERE id=?`, followInfo.FollowedSession).Scan(&followedid)
	if errq2 != nil {
		fmt.Println("Error get ID:", errq2)
		return
	}

	stmt, err := dataB.SocialDB.Prepare("UPDATE Notifications SET status = ? WHERE senderId = ? AND receiverId = ?")
	if err != nil {
		fmt.Println("Error change elem in db:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("refused", followerid, followedid)
	if err != nil {
		fmt.Println("Error change elem in db:", err)
		return
	}

	response := map[string]interface{}{
		"status": "follow declined successfuly",
	}
	json.NewEncoder(w).Encode(response)
}
