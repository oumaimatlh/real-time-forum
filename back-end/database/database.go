package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func DBinit() {
	var err error
	if DB, err = sql.Open("sqlite3", "./forum"); err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal(err)
	}
	TableCreation()
}
