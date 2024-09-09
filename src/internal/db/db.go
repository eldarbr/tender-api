package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dsn string) error {
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := DB.Ping(); err != nil {
		return err
	}
	return nil
}
