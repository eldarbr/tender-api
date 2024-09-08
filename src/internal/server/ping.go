package server

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		JSONError(w, map[string]string{"reason": "Method not allowed"}, 400)
		return
	}
	w.Write([]byte("ok"))
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	JSONError(w, map[string]string{"reason": "Not found"}, 400)
}
