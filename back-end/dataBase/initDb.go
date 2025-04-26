package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/sqlite3" 
    _ "github.com/golang-migrate/migrate/v4/source/file"      
    _ "github.com/mattn/go-sqlite3"
)

var socialDB *sql.DB

const dbPath = "dataBase/socialN.db"

func DbInit() {
	var err error
	socialDB, err = sql.Open("sqlite3", "./dataBase/socialN.db")
	if err != nil {
		fmt.Println("in 20 " , err)
		return
	} else {
		fmt.Println("db connection opened successfully ! ")
	}
	_, err = socialDB.Exec("PRAGMA foreign_keys = ON") ;  if err != nil {
		fmt.Println(err)
		return
	}
	err  = ApplyMigrations()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ApplyMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationsPath := "file://" + filepath.Join(wd, "dataBase", "migrations")

	// creating the migration
	m, err := migrate.New(migrationsPath, "sqlite3://"+dbPath)
	if err != nil {
		return err
	}
	// applying all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	
	log.Println("Migrations applied successfully!")
	return nil
}
