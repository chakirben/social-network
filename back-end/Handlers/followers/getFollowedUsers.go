package followers

import (
	"fmt"
	"strings"

	dataB "socialN/dataBase"
)

func GetFollowedUsers(loggedUserID int) []int {
	rows, err := dataB.SocialDB.Query(`
			SELECT followerId FROM Followers WHERE followedId=?
		`, loggedUserID)
	if err != nil {
		fmt.Println("Error :", err)
	}
	defer rows.Close()
	var followedUsersIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Error :", err)
		}
		followedUsersIds = append(followedUsersIds, id)
	}

	return followedUsersIds
}




func GetUnfollowedUsers(loggedUserID int, followedUsers []int) []string {
	var unfollowedUsersNames []string

	if len(followedUsers) == 0 {
		// If no one is followed, return all users
		rows, err := dataB.SocialDB.Query(`SELECT nickname FROM Users WHERE id!=?`, loggedUserID)
		if err != nil {
			fmt.Println("Error:", err)
			return unfollowedUsersNames
		}
		defer rows.Close()

		for rows.Next() {
			var nickname string
			if err := rows.Scan(&nickname); err != nil {
				fmt.Println("Error:", err)
			}
			unfollowedUsersNames = append(unfollowedUsersNames, nickname)
		}
		return unfollowedUsersNames
	}

	// Build the placeholders (?, ?, ?, ...)
	placeholders := make([]string, len(followedUsers)+1)
	args := make([]interface{}, len(followedUsers)+1)
	for i, id := range followedUsers {
		placeholders[i] = "?"
		args[i] = id
	}
	placeholders[len(placeholders)-1] = "?"
	args[len(args)-1] = loggedUserID

	query := fmt.Sprintf(
		`SELECT nickname FROM Users WHERE id NOT IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := dataB.SocialDB.Query(query, args...)
	if err != nil {
		fmt.Println("Error:", err)
		return unfollowedUsersNames
	}
	defer rows.Close()

	for rows.Next() {
		var nickname string
		if err := rows.Scan(&nickname); err != nil {
			fmt.Println("Error:", err)
		}
		unfollowedUsersNames = append(unfollowedUsersNames, nickname)
	}

	return unfollowedUsersNames
}
