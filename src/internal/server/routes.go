package server

import (
	"avito-back-test/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", handler.PingHandler).Methods(http.MethodGet)

	tenderHandler := handler.NewTenderHandler()
	r.HandleFunc("/api/tenders/new", tenderHandler.InsertNewTender).Methods(http.MethodPost)
	r.HandleFunc("/api/tenders/my", tenderHandler.GetMyTenders).Methods(http.MethodGet)
	r.HandleFunc("/api/tenders/{tenderId}/status", tenderHandler.UpdateTenderStatus).Methods(http.MethodPut)
	r.HandleFunc("/api/tenders/{tenderId}/status", tenderHandler.GetTenderStatus).Methods(http.MethodGet)
	r.HandleFunc("/api/tenders/{tenderId}/edit", tenderHandler.UpdateTender).Methods(http.MethodPatch)
	r.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender).Methods(http.MethodPut)
	r.HandleFunc("/api/tenders", tenderHandler.GetTenders).Methods(http.MethodGet)

	bidHandler := handler.NewBidHandler()
	r.HandleFunc("/api/bids/new", bidHandler.InsertNewBid).Methods(http.MethodPost)
	r.HandleFunc("/api/bids/my", bidHandler.GetMyBids).Methods(http.MethodGet)
	r.HandleFunc("/api/bids/{tenderId}/list", bidHandler.GetBidsByTender).Methods(http.MethodGet)
	r.HandleFunc("/api/bids/{bidId}/status", bidHandler.GetBidStatus).Methods(http.MethodGet)
	r.HandleFunc("/api/bids/{bidId}/status", bidHandler.UpdateBidStatus).Methods(http.MethodPut)
	r.HandleFunc("/api/bids/{bidId}/edit", bidHandler.UpdateBid).Methods(http.MethodPatch)

	// gorilla/mux:
	// Routes are tested in the order they were added to the router
	// If two routes match, the first one wins
	r.MethodNotAllowedHandler = handler.MethodNotAllowedHandler()
	r.NotFoundHandler = handler.NotFoundHandler()

	return r
}
