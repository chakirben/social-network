package profile

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	dataB "socialN/dataBase"
)

type ProfileOwner struct {
	ID string `json:"id"`
	Session string `json:"session"`
}

type FollowData struct {
	ID        int
	Nickname  string
	Firstname string
	Lastname  string
	Avatar    interface{}
	About     interface{}
}

func GetData(w http.ResponseWriter, r *http.Request) {
	var profile_owner ProfileOwner
	//decode the request into the struct
	err := json.NewDecoder(r.Body).Decode(&profile_owner)
	if err != nil {
		fmt.Println("Invalid Json:", err)
		return
	}

	//get account type of the user
	var accountType interface{}
	err = dataB.SocialDB.QueryRow("SELECT accountType FROM Users WHERE id=?", profile_owner.ID).Scan(&accountType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get the user id from the username
	user_id, _ := strconv.Atoi(profile_owner.ID)

	//get the user if of the loggedin
	var logged_user_id int
	err = dataB.SocialDB.QueryRow("SELECT userId FROM Sessions WHERE id=?", profile_owner.Session).Scan(&logged_user_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get personal data of profile owner
	var personal_data []interface{}
	var data FollowData
	err = dataB.SocialDB.QueryRow("SELECT nickname, firstName, lastName, avatar, about FROM Users WHERE id=?", user_id).Scan(&data.Nickname, &data.Firstname, &data.Lastname, &data.Avatar, &data.About)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	personal_data = append(personal_data, data)

	//count how much followers and followeds
	var followers_count int
	var followed_count int
	err = dataB.SocialDB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followedId=?", user_id).Scan(&followers_count)
	if err != nil {
		followers_count = 0
	}
	err = dataB.SocialDB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followerId=?", user_id).Scan(&followed_count)
	if err != nil {
		followed_count = 0
	}

	var followers []interface{}
	var followeds []interface{}
	posts := []map[string]interface{}{}
	if accountType == "public" {
		//get followers
		rows, errf := dataB.SocialDB.Query(`SELECT followerId FROM Followers WHERE followedId=?`, user_id)
		if errf != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var followerData FollowData
			var follower_id int
			err = rows.Scan(&follower_id)
			if err != nil {
				follower_id = 0
			}

			err = dataB.SocialDB.QueryRow("SELECT id, nickname, firstName, lastName, avatar, about FROM Users WHERE id=?", follower_id).Scan(&followerData.ID, &followerData.Nickname, &followerData.Firstname, &followerData.Lastname, &followerData.Avatar, &followerData.About)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			followers = append(followers, followerData)
		}

		//get followeds
		rowsd, errd := dataB.SocialDB.Query(`SELECT followedId FROM Followers WHERE followerId=?`, user_id)
		if errd != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rowsd.Close()
		for rowsd.Next() {
			var followedData FollowData
			var followed_id int
			err = rowsd.Scan(&followed_id)
			if err != nil {
				followed_id = 0
			}

			err = dataB.SocialDB.QueryRow("SELECT id, nickname, firstName, lastName, avatar, about FROM Users WHERE id=?", followed_id).Scan(&followedData.ID, &followedData.Nickname, &followedData.Firstname, &followedData.Lastname, &followedData.Avatar, &followedData.About)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			followeds = append(followeds, followedData)
		}


		//get posts
		query := `
			SELECT
				p.id,
				p.content,
				p.image,
				u.avatar,
				u.firstName,
				u.lastName,
				(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = 1) AS likeCount,
				(SELECT COUNT(*) FROM postReactions WHERE postId = p.id AND reactionType = -1) AS dislikeCount,
				(SELECT reactionType FROM postReactions WHERE postId = p.id AND userId = ?) AS userReaction,
				p.createdAt
			FROM Posts p
			JOIN Users u ON p.CreatorId = u.id
			WHERE p.creatorId = ?
			ORDER BY p.createdAt DESC
		`
		rows, err := dataB.SocialDB.Query(query, user_id, user_id)
		if err != nil {
			log.Println("Error fetching created posts:", err)
			http.Error(w, "Error fetching created posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var postID int
			var content, firstName, lastName string
			var image *string
			var avatar *string
			var likeCount, dislikeCount, userReaction int
			var createdAt time.Time

			err := rows.Scan(&postID, &content, &image, &avatar, &firstName, &lastName, &likeCount, &dislikeCount, &userReaction, &createdAt)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}

			post := map[string]interface{}{
				"id":            postID,
				"content":       content,
				"image":         image,
				"creator":       fmt.Sprintf("%s %s", firstName, lastName),
				"avatar":        avatar,
				"like_count":    likeCount,
				"dislike_count": dislikeCount,
				"user_reaction": userReaction,
				"created_at":    createdAt,
			}

			posts = append(posts, post)
		}
	}

	//check if the visitor is the profile ownwer
	var profile_status = "auther"
	if logged_user_id == user_id {
		profile_status = "mine"
	}


	//check if the logged use is already follow the profile owner
	var follow_status string
	if checkAlreadyFollow(logged_user_id, user_id) {
		follow_status = "Unfollow"
	} else {
		follow_status = "Follow"
	}

	//send response
	response := map[string]interface{}{
		"personal_data":   personal_data,
		"followers_count": followers_count,
		"followed_count":  followed_count,
		"followers_data":  followers,
		"followeds_data":  followeds,
		"posts":           posts,
		"follow_status": follow_status,
		"profile_status": profile_status,
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