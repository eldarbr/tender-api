package main

import (
	"avito-back-test/internal/config"
	"avito-back-test/internal/server"
	"os"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		// TODO: log the error
		os.Exit(1)
	}
	server := server.NewServer(config)
	server.ListenAndServe()
}
