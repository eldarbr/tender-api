package server

import (
	"avito-back-test/internal/config"
	"net/http"
	"time"
)

func NewServer(cfg *config.Config) *http.Server {
	router := newRouter()

	serv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 20,
	}

	return serv
}
