package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// The database package provides an interface between go code and a relational database.

var DB *sqlx.DB

func Init() {
	var err error

	DB, err = sqlx.Connect("sqlite3", "file:passwords.db")

	if err != nil {
		panic(err.Error())
	}

	// db.SetMaxOpenConns(0)
	// db.SetMaxIdleConns(200)
	// db.SetConnMaxLifetime(5 * time.Minute)
}
