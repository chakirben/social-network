package followers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataB "socialN/dataBase"
)

//this is just for public profile
//i will create for the followed profile if is private or public
//if public i will use this code
//if not, before add the follower i will send a response to front to show a notication of follower request
//if followed accept the follow request, i will receive the request come from front, then i will add follower in db

type FollowInfo struct {
	Follower string `json:"follower"`
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

	row2, errq2 := dataB.SocialDB.Query(`SELECT id FROM Users WHERE nickname=?`, followInfo.Follower)
	if errq2 != nil {
		fmt.Println("Error get ID:", errq2)
		return
	}
	defer row2.Close()
	for row2.Next() {
		errs2 := row2.Scan(&followerid)
		if errs2 != nil {
			fmt.Println("Error Scan:", errs2)
		}
	}


	var followed bool
	if !checkAlreadyFollow(followerid, followedid) {
		_, exec_err := dataB.SocialDB.Exec(`INSERT INTO Followers (followerId, followedId) VALUES (?,?)`, followerid, followedid)
		if exec_err != nil {
			fmt.Println("Error Insert into db:", exec_err)
			return
		}
		followed = true
	} else {
		_, exec_err := dataB.SocialDB.Exec(`DELETE FROM Followers WHERE followerId=? AND followedId=?`, followerid, followedid)
		if exec_err != nil {
			fmt.Println("Error delete in db:", exec_err)
			return
		}
		followed = false
	}
	


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":      "success",
		"followed": followed,
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