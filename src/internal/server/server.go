package server

import (
	"avito-back-test/internal/config"
	"net/http"
)

func NewServer(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", DefaultHandler)
	mux.HandleFunc("/api/ping", PingHandler)

	serv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: mux,
	}

	return serv
}
