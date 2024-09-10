package handler

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/repository"
	"avito-back-test/internal/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	// TODO: proper response in case of an empty list
	tenders, err := h.srv.GetTenders()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to fetch tenders", http.StatusInternalServerError)
		return
	}
	if tenders == nil {
		JSONResponse(w, []map[string]string{}, 200)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode tenders to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *TenderHandler) InsertNewTender(w http.ResponseWriter, r *http.Request) {
	var tenderRequest struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		ServiceType     string `json:"serviceType"`
		Status          string `json:"status"`
		OrganizationID  string `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&tenderRequest); err != nil {
		JSONResponse(w, map[string]string{"reason": "Invalid request payload"}, 400)
		return
	}

	// Convert OrganizationID to UUID
	orgID, err := uuid.Parse(tenderRequest.OrganizationID)
	if err != nil {
		JSONResponse(w, map[string]string{"reason": "Invalid organizationId format"}, 400)
		return
	}

	// Build the Tender model
	newTender := model.Tender{
		Name:            tenderRequest.Name,
		Description:     tenderRequest.Description,
		ServiceType:     tenderRequest.ServiceType,
		Status:          tenderRequest.Status,
		OrganizationID:  orgID,
		CreatorUsername: tenderRequest.CreatorUsername,
	}

	// Pass to the service
	err = h.srv.InsertNewTender(&newTender)
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": "the employee is not respnosible for the operation"}, 403)
		return
	}
	if err == repository.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": "no employee with set username"}, 401)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	JSONResponse(w, newTender, 200)
}
