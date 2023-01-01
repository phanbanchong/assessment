package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error

	command, err := os.ReadFile("sql/init.sql")
	if err != nil {
		log.Fatal("Read command to create database error", err)
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connection to database error", err)
	}

	_, err = db.Exec(string(command))
	if err != nil {
		log.Fatal("Unable to create table", err)
	} else {
		log.Printf("Executed: %s", command)
	}
}
