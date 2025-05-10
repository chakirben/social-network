package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataB "socialN/dataBase"
	followers "socialN/Handlers/followers"
)

type LoggedUser struct {
	Session string `json:"session"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var loggeduser LoggedUser

	err := json.NewDecoder(r.Body).Decode(&loggeduser)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	//get user id from session
	var user_id int
	err = dataB.SocialDB.QueryRow("SELECT userId FROM Sessions WHERE id=?", loggeduser.Session).Scan(&user_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}



	followed_users := followers.GetFollowedUsers(user_id)
	unfollowed_users := followers.GetUnfollowedUsers(user_id, followed_users)


	//send response
	response := map[string]interface{}{
		"unfollowed_users": unfollowed_users,
	}
	json.NewEncoder(w).Encode(response)
}
