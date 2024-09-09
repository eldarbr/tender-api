package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	DB = db
	return nil
}
