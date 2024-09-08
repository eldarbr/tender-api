package main

import (
	"avito-back-test/internal/config"
	"database/sql"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func disableSSL(connUrl string) (string, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return "", err
	}

	q := u.Query()

	q.Set("sslmode", "disable")

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	noSSLUrl, err := disableSSL(config.PostgresConnUrl)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", noSSLUrl)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		driver.Close()
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		m.Close()
		log.Fatal(err)
	}
	m.Close()
}
