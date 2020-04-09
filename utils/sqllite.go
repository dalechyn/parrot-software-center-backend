package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const dbPath = "test.db"
var db *sql.DB

func InitDB() {
	log.Info("Initializing SQLite3 database")

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("SQLite3 database not found in ", dbPath)
	}

	_, err = db.Exec("create table if not exists Ratings (package_name text primary key, package_rating real)")
	if err != nil {
		log.Debug("Initial table creation error", err)
	}
	log.Info("Initializing SQLite3 succeed")
}

func GetDB() *sql.DB {
	return db
}
