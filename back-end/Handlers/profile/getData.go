package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dataB "socialN/dataBase"
)

type ProfileOwner struct {
	ID string `json:"id"`
}

type FollowData struct {
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
	var accountType string
	err = dataB.SocialDB.QueryRow("SELECT accountType FROM Users WHERE id=?", profile_owner.ID).Scan(&accountType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get the user id from the username
	var user_id = profile_owner.ID


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

			err = dataB.SocialDB.QueryRow("SELECT nickname, firstName, lastName, avatar, about FROM Users WHERE id=?", follower_id).Scan(&followerData.Nickname, &followerData.Firstname, &followerData.Lastname, &followerData.Avatar, &followerData.About)
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

			err = dataB.SocialDB.QueryRow("SELECT nickname, firstName, lastName, avatar, about FROM Users WHERE id=?", followed_id).Scan(&followedData.Nickname, &followedData.Firstname, &followedData.Lastname, &followedData.Avatar, &followedData.About)
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
				p.title,
				p.content,
				p.image,
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()


		for rows.Next() {
			var postID int
			var title, content, firstName, lastName string
			var image sql.NullString
			var likeCount, dislikeCount, userReaction sql.NullInt32
			var createdAt time.Time

			err := rows.Scan(&postID, &title, &content, &image, &firstName, &lastName, &likeCount, &dislikeCount, &userReaction, &createdAt)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}

			post := map[string]interface{}{
				"id":           postID,
				"title":        title,
				"content":      content,
				"image":        image,
				"creator":      fmt.Sprintf("%s %s", firstName, lastName),
				"likeCount":    likeCount.Int32,
				"dislikeCount": dislikeCount.Int32,
				"userReaction": userReaction.Int32,
				"createdAt":    createdAt,
			}

			posts = append(posts, post)
		}
	}

	//send response
	response := map[string]interface{}{
		"personal_data": personal_data,
		"followers_count": followers_count,
		"followed_count":  followed_count,
		"followers_data":  followers,
		"followeds_data":  followeds,
		"posts": posts,
	}
	json.NewEncoder(w).Encode(response)
}
