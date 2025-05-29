package followers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"socialN/Handlers/auth"
	dataB "socialN/dataBase"
)

// there is two function: one function if someone follow auther it handle if the profile is public or private
// and auther function handle if the followed accept follow request of follower
type FollowInfo struct {
	Follower_session string `json:"follower_session"`
	Followed_id      string `json:"followed_id"`
}

func HandleFollow(w http.ResponseWriter, r *http.Request) {
	var followInfo FollowInfo

	var followerid int
	err := json.NewDecoder(r.Body).Decode(&followInfo)
	if err != nil {
		/*fmt.Println("Invalid Json:", err)
		return*/
		receiverIdStr := r.URL.Query().Get("id")
		if receiverIdStr == "" {
			http.Error(w, "Invalid receiver id", http.StatusBadRequest)
			return
		}

		senderId, err := auth.ValidateSession(r, dataB.SocialDB)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		followInfo.Followed_id = receiverIdStr
		followerid = senderId
	}

	followedid, err := strconv.Atoi(followInfo.Followed_id)
	//var followerid int
	var followedtype string
	errq := dataB.SocialDB.QueryRow(`SELECT accountType FROM Users WHERE id=?`, followedid).Scan(&followedtype)
	if errq != nil || err != nil {
		fmt.Println("Error get ID:", errq)
		return
	}

	if followInfo.Follower_session != "" {
		row2, errq2 := dataB.SocialDB.Query(`SELECT userId FROM Sessions WHERE id=?`, followInfo.Follower_session)
		if errq2 != nil {
			fmt.Println("Error get ID:", errq2)
			return
		}
		defer row2.Close()
		for row2.Next() {
			errs2 := row2.Scan(&followerid)
			if errs2 != nil {
				fmt.Println("Error Scan:", errs2)
				return
			}
		}
	}
	

	var response map[string]interface{}

	if followedtype == "public" {
		var followed bool
		var followedText string
		if !checkAlreadyFollow(followerid, followedid) {
			_, exec_err := dataB.SocialDB.Exec(`INSERT INTO Followers (followerId, followedId) VALUES (?,?)`, followerid, followedid)
			if exec_err != nil {
				fmt.Println("Error Insert into db:", exec_err)
				return
			}
			followed = true
			followedText = "Unfollow"
		} else {
			_, exec_err := dataB.SocialDB.Exec(`DELETE FROM Followers WHERE followerId=? AND followedId=?`, followerid, followedid)
			if exec_err != nil {
				fmt.Println("Error delete in db:", exec_err)
				return
			}
			followed = false
			followedText = "Follow"
		}

		var followers_count int
		err = dataB.SocialDB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followedId=?", followedid).Scan(&followers_count)
		if err != nil {
			followers_count = 0
		}

		response = map[string]interface{}{
			"status":          followedText,
			"followed":        followed,
			"followers_count": followers_count,
		}
	} else if followedtype == "private" {

		var followedText = "waiting for followed accept"
		if checkAlreadyFollow(followerid, followedid) {
			_, exec_err := dataB.SocialDB.Exec(`DELETE FROM Followers WHERE followerId=? AND followedId=?`, followerid, followedid)
			if exec_err != nil {
				fmt.Println("Error delete in db:", exec_err)
				return
			}
			followedText = "Follow"
		} else {
			//insert into notifications table
			_, err = dataB.SocialDB.Exec(`
				INSERT INTO Notifications (senderId, receiverId, type, status, notificationDate)
				VALUES (?, ?, 'follow_request', 'pending', ?)
			`, followerid, followedid, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Println("Error inserting notification:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		
		
		var followers_count int
		err = dataB.SocialDB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followedId=?", followedid).Scan(&followers_count)
		if err != nil {
			followers_count = 0
		}
		

		response = map[string]interface{}{
			"status":          followedText,
			"followed":        followedid,
			"follower":        followerid,
			"followers_count": followers_count,
		}
	}

	json.NewEncoder(w).Encode(response)
}

func checkAlreadyFollow(followerID, followedID int) bool {
	rows, err := dataB.SocialDB.Query(`SELECT 1 FROM Followers WHERE followerId = ? AND followedId = ?`, followerID, followedID)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}


type acceptFollow struct {
	FollowerID int `json:"follower_id"`
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