package followers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataB "socialN/dataBase"
)

//there is two function: one function if someone follow auther it handle if the profile is public or private
//and auther function handle if the followed accept follow request of follower

type FollowInfo struct {
	Follower_session string `json:"follower_session"`
	Followed string `json:"followed"`
}


func HandleFollow(w http.ResponseWriter, r *http.Request) {
	var followInfo FollowInfo

	err := json.NewDecoder(r.Body).Decode(&followInfo)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	var followerid int
	var followedid int
	var followedtype string
	row, errq := dataB.SocialDB.Query(`SELECT id, accountType FROM Users WHERE nickname=?`, followInfo.Followed)
	if errq != nil {
		fmt.Println("Error get ID:", errq)
		return
	}
	defer row.Close()
	for row.Next() {
		errs := row.Scan(&followedid, &followedtype)
		if errs != nil {
			fmt.Println("Error Scan:", err)
			return
		}
	}

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
			followedText = "followed successfuly"
		} else {
			_, exec_err := dataB.SocialDB.Exec(`DELETE FROM Followers WHERE followerId=? AND followedId=?`, followerid, followedid)
			if exec_err != nil {
				fmt.Println("Error delete in db:", exec_err)
				return
			}
			followed = false
			followedText = "unfollowed successfuly"
		}


		response = map[string]interface{}{
			"status": followedText,
			"followed": followed,
		}
	} else if followedtype == "private" {
		response = map[string]interface{}{
			"status":      "waiting for followed accept",
			"followed": followInfo.Followed,
			"follower": followInfo.Follower_session,
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




func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	var followInfo FollowInfo

	err := json.NewDecoder(r.Body).Decode(&followInfo)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	var followerid int
	var followedid int
	row, errq := dataB.SocialDB.Query(`SELECT id FROM Users WHERE nickname=?`, followInfo.Followed)
	if errq != nil {
		fmt.Println("Error get ID:", errq)
		return
	}
	defer row.Close()
	for row.Next() {
		errs := row.Scan(&followedid)
		if errs != nil {
			fmt.Println("Error Scan:", err)
			return
		}
	}

	row2, errq2 := dataB.SocialDB.Query(`SELECT userId FROM Users WHERE id=?`, followInfo.Follower_session)
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

	_, exec_err := dataB.SocialDB.Exec(`INSERT INTO Followers (followerId, followedId) VALUES (?,?)`, followerid, followedid)
	if exec_err != nil {
		fmt.Println("Error Insert into db:", exec_err)
		return
	}

	response := map[string]interface{}{
		"status":      "followed successfuly",
	}
	json.NewEncoder(w).Encode(response)
}