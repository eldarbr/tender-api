package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", PingHandler).Methods(http.MethodGet)

	// gorilla/mux:
	// Routes are tested in the order they were added to the router
	// If two routes match, the first one wins
	r.MethodNotAllowedHandler = MethodNotAllowedHandler()
	r.NotFoundHandler = NotFoundHandler()

	return r
}
