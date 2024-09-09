package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.MethodNotAllowedHandler = MethodNotAllowedHandler()
	r.PathPrefix("/").HandlerFunc(DefaultHandler)
	r.HandleFunc("/api/ping", PingHandler).Methods(http.MethodGet)

	return r
}
