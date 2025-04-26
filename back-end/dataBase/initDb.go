package database

import (
	"database/sql"
	"fmt"
)

var ForumDB *sql.DB

func DbInit() {
	var err error
	ForumDB, err = sql.Open("sqlite3", "./dataBase/forum.db")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("success")
	}
	CreateTables(ForumDB)
}
