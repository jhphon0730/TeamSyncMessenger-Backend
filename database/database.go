package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Database Connect Success.")
	return db
}
