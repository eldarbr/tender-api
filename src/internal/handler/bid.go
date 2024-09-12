package handler

import (
	"avito-back-test/internal/model"
	"avito-back-test/internal/service"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BidHandler struct {
	srv *service.BidService
}

func NewBidHandler() *BidHandler {
	srv := service.NewBidService()
	return &BidHandler{
		srv: srv,
	}
}

func (h *BidHandler) InsertNewBid(w http.ResponseWriter, r *http.Request) {
	var bidRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		TenderID    string `json:"tenderId"`
		AuthorType  string `json:"authorType"`
		AuthorID    string `json:"authorId"`
	}

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
		JSONResponse(w, map[string]string{"reason": "Invalid request payload"}, 400)
		return
	}

	// Convert TenderID to UUID
	tendID, err1 := uuid.Parse(bidRequest.TenderID)
	authorId, err2 := uuid.Parse(bidRequest.AuthorID)
	if err1 != nil || err2 != nil {
		JSONResponse(w, map[string]string{"reason": "Invalid organizationId format"}, 400)
		return
	}

	// Build the Bid model
	newBid := model.Bid{
		Name:        bidRequest.Name,
		Description: bidRequest.Description,
		TenderID:    tendID,
		AuthorType:  bidRequest.AuthorType,
		AuthorID:    authorId,
	}

	// Pass to the service
	err := h.srv.InsertNewBid(&newBid)
	if err == service.ErrNoEmployee || err == service.ErrNoOrganization {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
	JSONResponse(w, newBid, 200)
}

func (h *BidHandler) GetMyBids(w http.ResponseWriter, r *http.Request) {
	var (
		bids []model.Bid
		err  error

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

	bids, err = h.srv.GetUserBids(username[0], limit, offset)
	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	if bids == nil {
		JSONResponse(w, []map[string]string{}, 200)
		return
	}
	if err := json.NewEncoder(w).Encode(bids); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
}

func (h *BidHandler) GetBidsByTender(w http.ResponseWriter, r *http.Request) {
	var (
		bids []model.Bid
		err  error

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
	requestVars := mux.Vars(r)
	tenderID, err := uuid.Parse(requestVars["tenderId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}

	bids, err = h.srv.GetBidsByTender(tenderID, username[0], limit, offset)

	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err == service.ErrNoTender {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	if bids == nil {
		JSONResponse(w, []map[string]string{}, 200)
		return
	}
	if err := json.NewEncoder(w).Encode(bids); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 500)
		return
	}
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	username, ok := queryValues["username"]
	if !ok {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}
	requestVars := mux.Vars(r)
	bidID, err := uuid.Parse(requestVars["bidId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}

	status, err := h.srv.GetBidStatus(bidID, username[0])

	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err == service.ErrNoBid {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, status, 200)
}

func (h *BidHandler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	username, ok := queryValues["username"]
	if !ok {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}
	status, ok := queryValues["username"]
	if !ok {
		JSONResponse(w, map[string]string{"reason": "username is required"}, 400)
		return
	}
	requestVars := mux.Vars(r)
	bidID, err := uuid.Parse(requestVars["bidId"])
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	bid := model.Bid{
		ID:     bidID,
		Status: status[0],
	}
	err = h.srv.UpdateBidStatus(&bid, username[0])

	if err == service.ErrNoEmployee {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 401)
		return
	}
	if err == service.ErrNoBid {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 404)
		return
	}
	if err == service.ErrNotResponsible {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 403)
		return
	}
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, bid, 200)
}

func (h *BidHandler) UpdateBid(w http.ResponseWriter, r *http.Request) {
	var (
		bidUpdate model.BidUpdate
		username  string
	)
	if err := json.NewDecoder(r.Body).Decode(&bidUpdate); err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidId"])
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

	updatedBid, err := h.srv.PatchBid(bidID, username, &bidUpdate)
	if err == service.ErrNoBid {
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
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, *updatedBid, 200)
}

func (h *BidHandler) RollbackBid(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		version  int
	)
	vars := mux.Vars(r)
	bidID, err := uuid.Parse(vars["bidId"])
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

	updatedBid, err := h.srv.RollbackBid(bidID, username, version)
	if err == service.ErrNoBid {
		JSONResponse(w, map[string]string{"reason": "no bid with specified version"}, 404)
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
	if err != nil {
		JSONResponse(w, map[string]string{"reason": err.Error()}, 400)
		return
	}
	JSONResponse(w, *updatedBid, 200)
}
