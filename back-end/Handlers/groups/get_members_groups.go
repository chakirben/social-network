package groups

import (
	dataB "socialN/dataBase"
)

func GetMembersGroups(idgroup int) (int, error) {
	query := `
	SELECT count(*)
	FROM GroupsMembers
	WHERE groupId = ?
	`
	var cont int
	err := dataB.SocialDB.QueryRow(query, idgroup).Scan(&cont)
	if err != nil {
		return -1, err
	}
	return cont, nil
}
