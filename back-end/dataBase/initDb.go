package dataB

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

var SocialDB *sql.DB

const dbPath = "dataBase/socialN.db"

func DbInit() {
	var err error
	SocialDB, err = sql.Open("sqlite3", "../dataBase/socialN.db")
	if err != nil {
		fmt.Println("in 20 ", err)
		return
	} else {
		fmt.Println("db connection opened successfully ! ")
	}
	_, err = SocialDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ApplyMigrations()
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
	parentDir := filepath.Dir(wd)
	migrationsDir := filepath.Join(parentDir, "dataBase", "migrations")
	migrationsPath := "file://" + filepath.ToSlash(migrationsDir)
	absDbPath := filepath.Join(parentDir, "dataBase", "socialN.db")
	dbUri := "sqlite3://" + absDbPath

	m, err := migrate.New(migrationsPath, dbUri)
	if err != nil {
		fmt.Println("Migration create error:", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully!")
	return nil
}
