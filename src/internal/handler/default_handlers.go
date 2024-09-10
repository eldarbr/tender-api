package handler

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", http.MethodGet)
		JSONResponse(w, map[string]string{"reason": "Method not allowed"}, 405)
	})
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, map[string]string{"reason": "Not found"}, 404)
	})
}
