package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"parrot-software-center-backend/models"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const dbPath = "test.db"
var db *sql.DB

func InitUserTable() {
	_, err := db.Exec("create table if not exists users (id integer primary key autoincrement, " +
		"email text not null, username text not null, password text not null)")
	if err != nil {
		log.Fatal("Initial table creation error", err)
	}
}

func InitRatingsTable() {
	_, err := db.Exec("create table if not exists ratings (user_id integer primary key not null, " +
		"name text not null, author text not null, rating integer not null, commentary text)")
	if err != nil {
		log.Fatal("Initial table creation error", err)
	}
}

func InitDB() {
	log.Info("Initializing database")

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("SQLite3 database not found in ", dbPath)
	}

	InitUserTable()
	InitRatingsTable()
	log.Info("Initializing database succeed")
}

func GetDB() *sql.DB {
	return db
}

func GetIDFromToken(tokenStr string) (int, error) {
	hmacSecret := []byte(GetSecret())
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSecret, nil
	})

	if err != nil {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*models.Claims);
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Check if user exists
	id := -1
	row := db.QueryRow("select id from users where username = $1", claims.Username)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}