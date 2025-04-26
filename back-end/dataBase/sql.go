package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTables(db *sql.DB) error {
	statement, err := os.ReadFile("./dataBase/schema.sql")
	if err != nil {
		fmt.Println("err readinfile", err)
		return err
	}
	_, err = db.Exec(string(statement))
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	return nil
}
