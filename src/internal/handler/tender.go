package handler

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func parseQueryLimitOffset(query *url.Values) (int, int, error) {
	var (
		limit  = 5
		offset = 0
		err    error
	)
	sLimit, ok := (*query)["limit"]
	if ok {
		limit, err = strconv.Atoi(sLimit[0])
		if err == nil && limit < 1 {
			return 0, 0, errors.New("limit has to be positive")
		}
	}

	if err != nil {
		return 0, 0, err
	}

	sOffset, ok := (*query)["offset"]
	if ok {
		offset, err = strconv.Atoi(sOffset[0])
		if err == nil && limit < 1 {
			return 0, 0, errors.New("offset has to be non-negative")
		}
	}

	return limit, offset, nil
}

func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	var (
		tenders []model.Tender
		err     error

		// query parameters
		limit, offset int
	)

	queryValues := r.URL.Query()
	limit, offset, err = parseQueryLimitOffset(&queryValues)
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}

	if serviceType, ok := queryValues["service_type"]; ok {
		tenders, err = h.srv.GetTendersOfService(serviceType[0], limit, offset)
	} else {
		tenders, err = h.srv.GetTenders(limit, offset)
	}

	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	if tenders == nil {
		JSONResponse(w, []map[string]string{}, 200)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
}

func (h *TenderHandler) InsertNewTender(w http.ResponseWriter, r *http.Request) {
	var tenderRequest struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		ServiceType     string `json:"serviceType"`
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
		Name:           tenderRequest.Name,
		Description:    tenderRequest.Description,
		ServiceType:    tenderRequest.ServiceType,
		OrganizationID: orgID,
	}

	// Pass to the service
	err = h.srv.InsertNewTender(&newTender, tenderRequest.CreatorUsername)
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": "the employee is not respnosible for the organization"}, 403)
		return
	}
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": "no employee with set username"}, 401)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
	JSONResponse(w, newTender, 200)
}

func (h *TenderHandler) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	var (
		tenders []model.Tender
		err     error

		// query parameters
		limit, offset int
		username      []string
	)

	queryValues := r.URL.Query()
	limit, offset, err = parseQueryLimitOffset(&queryValues)
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	username, ok := queryValues["username"]
	if !ok {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}

	tenders, err = h.srv.GetUserTenders(username[0], limit, offset)

	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	if tenders == nil {
		JSONResponse(w, []map[string]string{}, 200)
		return
	}
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
}

func (h *TenderHandler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	var err error

	requestVars := mux.Vars(r)

	r.ParseForm()
	if !r.Form.Has("status") || !r.Form.Has("username") {
		JSONResponse(w, map[string]string{"reason": "status, username are required"}, 400)
		return
	}
	var tender model.Tender
	tender.ID, err = uuid.Parse(requestVars["tenderId"])
	username := r.Form.Get("username")
	tender.Status = r.Form.Get("status")
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}

	err = h.srv.UpdateTenderStatus(&tender, username)

	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible || err == service.ErrTenderClosed {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	JSONResponse(w, tender, 200)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		tenderID uuid.UUID
		username *string
	)
	r.ParseForm()
	requestVars := mux.Vars(r)
	tenderID, err = uuid.Parse(requestVars["tenderId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	if r.Form.Has("username") {
		username = new(string)
		*username = r.Form.Get("username")
	}

	tenderStatus, err := h.srv.GetTenderStatus(tenderID, username)

	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": "not authorized"}, 403)
		return
	}
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	JSONResponse(w, tenderStatus, 200)
}

func (h *TenderHandler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	var (
		tenderUpdate model.TenderUpdate
		username     string
	)
	if err := json.NewDecoder(r.Body).Decode(&tenderUpdate); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	vars := mux.Vars(r)
	tenderID, err := uuid.Parse(vars["tenderId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	r.ParseForm()
	if r.Form.Has("username") {
		username = r.Form.Get("username")
	} else {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}

	updatedTender, err := h.srv.PatchTender(tenderID, username, &tenderUpdate)
	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible || err == service.ErrTenderClosed {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, *updatedTender, 200)
}

func (h *TenderHandler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		version  int
	)
	vars := mux.Vars(r)
	tenderID, err := uuid.Parse(vars["tenderId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	versionS, ok := vars["version"]
	if !ok {
		JSONResponse(w, map[string]string{"reason": "version is required"}, 400)
		return
	}
	version, err = strconv.Atoi(versionS)
	if err != nil {
		JSONResponse(w, map[string]string{"reason": "invalid version"}, 400)
		return
	}
	r.ParseForm()
	if r.Form.Has("username") {
		username = r.Form.Get("username")
	} else {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}

	updatedTender, err := h.srv.RollbackTender(tenderID, username, version)
	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": "no tender with specified version"}, 404)
		return
	}
	if err == service.ErrNotResponsible || err == service.ErrTenderClosed {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, *updatedTender, 200)
}
