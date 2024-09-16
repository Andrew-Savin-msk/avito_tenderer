package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	"github.com/gorilla/mux"
)

// Bid endpoints

func (s *server) handleCreateBid() http.HandlerFunc {
	type request struct {
		Name       string `json:"name"`
		Descr      string `json:"description"`
		TenderId   string `json:"tenderId"`
		AuthorType string `json:"authorType"`
		AuthorId   string `json:"authorId"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get into model struct
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		b := &models.Bid{
			Name:        req.Name,
			Description: req.Descr,
			TenderId:    req.TenderId,
			AuthorType:  req.AuthorType,
			AuthorId:    req.AuthorId,
		}
		// validate
		err = b.Validate()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		// BidsServ.Create()
		data, err := s.BidsServ.Create(b)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSuchTender) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleGetUsersBids() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse querry: limit, offset, username
		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		var err error
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusUnauthorized, ErrInvalidQuerryParams)
			return
		}

		// BidsServ.GetByName()
		data, err := s.BidsServ.GetByName(limit, offset, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce [data, data, data]
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleGetTendersBids() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: tenderId
		tenderId := mux.Vars(r)["tenderId"]
		if tenderId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusUnauthorized, ErrInvalidQuerryParams)
			return
		}
		// parse querry: limit, offset
		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		var err error
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}
		// BidsServ.GetTendersBids()
		data, err := s.BidsServ.GetTenderBids(limit, offset, tenderId, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSucnResource) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// respoce [data, data, data]
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleInterractBidStatus() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// parse path: bidId
			bidId := mux.Vars(r)["bidId"]
			if bidId == "" {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}

			// parse querry: username
			username := r.URL.Query().Get("username")
			if username == "" {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidQuerryParams)
				return
			}
			// BidsServ.GetStat()
			status, err := s.BidsServ.GetStat(bidId, username)
			if err != nil {
				if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
					s.DeadOnError(err)
					s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
					return
				}
				if errors.Is(err, services.ErrNoSuchUser) {
					s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
					return
				}
				if errors.Is(err, services.ErrNoPermitions) {
					s.error(w, r, http.StatusForbidden, err)
					return
				}
				if errors.Is(err, services.ErrNoSuchBid) {
					s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
					return
				}
				s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
				return
			}
			// response status
			s.respond(w, r, http.StatusOK, status)
			return
		} else if r.Method == http.MethodPut {
			// parse path: bidId
			bidId := mux.Vars(r)["bidId"]
			if bidId == "" {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// parse querry: status, username
			status := r.URL.Query().Get("status")
			if status == "" || (status != "Created" && status != "Published" && status != "Canceled") {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}

			username := r.URL.Query().Get("username")
			if username == "" {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// BidsServ.ChangeStat()
			data, err := s.BidsServ.ChangeStat(bidId, status, username)
			if err != nil {
				if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
					s.DeadOnError(err)
					s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
					return
				}
				if errors.Is(err, services.ErrNoSuchUser) {
					s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
					return
				}
				if errors.Is(err, services.ErrNoPermitions) {
					s.error(w, r, http.StatusForbidden, err)
					return
				}
				if errors.Is(err, services.ErrNoSuchBid) {
					s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
					return
				}
				s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
				return
			}
			// response data
			s.respond(w, r, http.StatusOK, data)
			return
		}
		s.error(w, r, http.StatusMethodNotAllowed, ErrUnsupportedMethod)
	})
}

func (s *server) handleEditBid() http.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Descr string `json:"description"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse body data
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		b := &models.Bid{
			Name:        req.Name,
			Description: req.Descr,
		}

		// parse path: bidId
		bidId := mux.Vars(r)["bidId"]
		if bidId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: username
		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// BidsServ.Edit()
		data, err := s.BidsServ.Edit(b, bidId, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSuchBid) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleSumbitBidDecision() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: bidId
		bidId := mux.Vars(r)["bidId"]
		if bidId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: descision, username
		decision := r.URL.Query().Get("decision ")
		if decision == "" || (decision != "Approved" && decision != "Rejected") {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// BidServ.Sumbit()
		data, err := s.BidsServ.Sumbit(bidId, decision, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSuchBid) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleBidFeedback() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: bidId
		bidId := mux.Vars(r)["bidId"]
		if bidId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: bidFeedback, username
		bidFeedback := r.URL.Query().Get("bidFeedback")
		if bidFeedback == "" || len(bidFeedback) > 1000 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// BidServ.AddFeedback()
		data, err := s.BidsServ.AddFeedback(bidId, bidFeedback, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSucnResource) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleRollbackBid() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: bidId, version
		bidId := mux.Vars(r)["bidId"]
		if bidId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		version, err := strconv.ParseInt(mux.Vars(r)["version"], 10, 32)
		if err != nil || version < 0 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: username
		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// BidServ.Rollback()
		data, err := s.BidsServ.Rollback(bidId, version, username)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSuchBid) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleGetTenderBidsReviews() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: tenderId
		tenderId := mux.Vars(r)["tenderId"]
		if tenderId == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: authorUsername, requesterUsername, limit, offset
		authorUsername := r.URL.Query().Get("authorUsername")
		if authorUsername == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		requesterUsername := r.URL.Query().Get("requesterUsername")
		if requesterUsername == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		var err error
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}
		// BidServ.GetReviews()
		data, err := s.BidsServ.GetReviews(tenderId, authorUsername, requesterUsername, limit, offset)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			if errors.Is(err, services.ErrNoSuchUser) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
				return
			}
			if errors.Is(err, services.ErrNoPermitions) {
				s.error(w, r, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, services.ErrNoSucnResource) {
				s.error(w, r, http.StatusNotFound, ErrNoSuchResorce)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce [data, data, data]
		s.respond(w, r, http.StatusOK, data)
	})
}
