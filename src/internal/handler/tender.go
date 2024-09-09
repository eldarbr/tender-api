package handler

import (
	"avito-back-test/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type TenderHandler struct {
	srv *service.TenderService
}

func NewTenderHandler() *TenderHandler {
	srv := service.NewTenderService()
	return &TenderHandler{
		srv: srv,
	}
}

func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	tenders, err := h.srv.GetTenders()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to fetch tenders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		log.Println(err)

		http.Error(w, "Failed to encode tenders to JSON", http.StatusInternalServerError)
		return
	}
}
