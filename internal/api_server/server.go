package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// sessionName        = "Authorization"
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type availability struct {
	is     bool
	reason error
}

type server struct {
	router *mux.Router
	logger *logrus.Logger

	TendersServ services.Tenders
	BidsServ    services.Bids

	available availability
}

func newServer(logger *logrus.Logger, TendersServ services.Tenders, BidsServ services.Bids) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logger,

		TendersServ: TendersServ,
		BidsServ:    BidsServ,

		available: availability{
			is: true,
		},
	}

	srv.configureRouter()

	return srv
}

func (s *server) DeadOnError(err error) {
	s.available.is = false
	s.available.reason = err
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router = s.router.PathPrefix("/api").Subrouter()

	s.router.Use(s.setRequestID)
	s.router.Use(s.deadChecker)
	s.router.Use(s.logRequest)
	s.router.Use(s.recoverPanic)

	// Tenders endpoints
	s.router.HandleFunc("/ping", s.handlePing()).Methods("GET")
	s.router.HandleFunc("/tenders", s.handleGetTendersList()).Methods("GET")
	s.router.HandleFunc("/tenders/new", s.handleCreateTender()).Methods("POST")
	s.router.HandleFunc("/tenders/my", s.handleGetUsersTenders()).Methods("GET")
	s.router.HandleFunc("/tenders/{tenderId}/status", s.handleInterractTenderStatus()).Methods("GET", "PUT")
	s.router.HandleFunc("/tenders/{tenderId}/edit", s.handleEditTender()).Methods("PATCH")
	s.router.HandleFunc("/tenders/{tenderId}/rollback/{version}", s.handleRollbackTender()).Methods("PUT")
	// Bids endpoints
	s.router.HandleFunc("/bids/new", s.handleCreateBid()).Methods("POST")
	s.router.HandleFunc("/bids/my", s.handleGetUsersBids()).Methods("GET")
	s.router.HandleFunc("/bids/{tenderId}/list", s.handleGetTendersBids()).Methods("GET")
	s.router.HandleFunc("/bids/{bidId}/status", s.handleInterractBidStatus()).Methods("GET", "PUT")
	s.router.HandleFunc("/bids/{bidId}/edit", s.handleEditBid()).Methods("PATCH")
	s.router.HandleFunc("/bids/{bidId}/sumbit_decision", s.handleSumbitBidDecision()).Methods("PUT")
	s.router.HandleFunc("/bids/{bidId}/feedback", s.handleBidFeedback()).Methods("PUT")
	s.router.HandleFunc("/bids/{bidId}/rallback/{version}", s.handleRollbackBid()).Methods("PUT")
	s.router.HandleFunc("/bids/{tenderId}/reviews", s.handleGetTenderBidsReviews()).Methods("GET")
}

// Func for making call of respond func with Error pattern
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"reason": err.Error()})
}

// Universal func for sending any type of respond (Error, Responde, etc.)
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
