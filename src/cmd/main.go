package main

import (
	"avito-back-test/internal/config"
	"avito-back-test/internal/server"
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
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
