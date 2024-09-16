package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/domain/models"
	"github.com/avito-testirovanie-na-backend-1270/cnrprod1725729417-team-78417/zadanie-6105/internal/services"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

// Tender endpoints

func (s *server) handleGetTendersList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse querry: limit, offset, service_type
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

		serviceTypes := r.URL.Query()["service_type"]
		if len(serviceTypes) != 0 {
			err = validation.Validate(serviceTypes,
				validation.Each(validation.In("Construction", "Delivery", "Manufacture")),
			)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		// Tenders.List()

		data, err := s.TendersServ.List(limit, offset, serviceTypes)
		if err != nil {
			if errors.Is(err, services.ErrServiceDatabaseDisconnected) {
				s.DeadOnError(err)
				s.error(w, r, http.StatusInternalServerError, ErrServiceUnavailable)
				return
			}
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}

		// responce [data, data, data]
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleCreateTender() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Descr    string `json:"description"`
		ServType string `json:"serviceType"`
		OrgId    string `json:"organizationId"`
		Username string `json:"creatorUsername"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		// get into model struct
		t := &models.Tender{
			Name:        req.Name,
			Description: req.Descr,
			ServType:    req.ServType,
		}

		// validate
		err = t.Validate()
		if err != nil || len(req.OrgId) > 100 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		ru := &models.Responsible{
			OrgId:    req.OrgId,
			Username: req.Username,
		}

		// TendersServ.Create()
		data, err := s.TendersServ.Create(t, ru)
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

		// responce data
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleGetUsersTenders() http.HandlerFunc {
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
		// TenserServ.GetByName()
		data, err := s.TendersServ.GetByName(limit, offset, username)
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
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}
		// responce [data, data, data]
		s.respond(w, r, http.StatusOK, data)
	})
}

func (s *server) handleInterractTenderStatus() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// parse path: tenderId
			tenderId := mux.Vars(r)["tenderId"]
			if tenderId == "" || len(tenderId) > 100 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// parse querry: username
			username := r.URL.Query().Get("username")
			if username == "" {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// TenderServ.GetStat()
			data, err := s.TendersServ.GetStat(tenderId, username)
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
			// response status
			s.respond(w, r, http.StatusOK, data)
		} else if r.Method == http.MethodPut {
			// parse path: tenderId
			tenderId := mux.Vars(r)["tenderId"]
			if tenderId == "" || len(tenderId) > 100 {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// parse querry: status, username
			status := r.URL.Query().Get("status")
			if status == "" || (status != "Created" && status != "Published" && status != "Closed") {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}

			username := r.URL.Query().Get("username")
			if username == "" {
				s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
			// TenderServ.ChangeStat()
			data, err := s.TendersServ.ChangeStat(tenderId, status, username)
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
			// response data
			s.respond(w, r, http.StatusOK, data)
		} else {
			s.error(w, r, http.StatusMethodNotAllowed, ErrUnsupportedMethod)
		}
	})
}

func (s *server) handleEditTender() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Descr    string `json:"description"`
		ServType string `json:"serviceType"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		// get into model struct
		t := &models.Tender{
			Name:        req.Name,
			Description: req.Descr,
			ServType:    req.ServType,
		}

		// parse path: tenderId
		tenderId := mux.Vars(r)["tenderId"]
		if tenderId == "" || len(tenderId) > 100 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: username
		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}

		err = t.ValidateEdition()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidRequestBody)
			return
		}

		// TenderServ.Edit()
		data, err := s.TendersServ.Edit(t, tenderId, username)
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

func (s *server) handleRollbackTender() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse path: tenderId, version
		tenderId := mux.Vars(r)["tenderId"]
		if tenderId == "" || len(tenderId) > 100 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		version, err := strconv.ParseInt(mux.Vars(r)["version"], 10, 32)
		if err != nil || version < 1 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// parse querry: username
		username := r.URL.Query().Get("username")
		if username == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
			return
		}
		// Tender.Rollback()
		data, err := s.TendersServ.Rollback(tenderId, username, version)
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
