package main

import (
	"avito-back-test/internal/config"
	"avito-back-test/internal/db"
	"avito-back-test/internal/server"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

func initDB(dsn string) error {
	// Channel to signal completion of the function
	done := make(chan bool, 1)
	var err error
	// Run the function in a goroutine
	go func() {
		err = db.InitDB(dsn)
		done <- true
	}()

	// Set a timeout duration
	timeout := time.After(10 * time.Second)

	// Wait for either the function to complete or the timeout
	select {
	case <-done:
		return err
	case <-timeout:
		return errors.New("db connection timed out")
	}
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := initDB(config.PostgresConnUrl); err != nil {
		log.Fatal(err)
	} else {
		log.Println("db init complete")
	}
	defer db.DB.Close()

	server := server.NewServer(config)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	context, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	server.Shutdown(context)

	log.Println("shutting down")
	os.Exit(0)
}
