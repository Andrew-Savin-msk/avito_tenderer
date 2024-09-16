package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// setRequestID sets UUID for every request
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// logRequest loggin any request and it's responce
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remout_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed in %v %s in %v",
			rw.code,
			http.StatusText((rw.code)),
			time.Now().Sub(start),
		)
	})
}

func (s *server) deadChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.available.is {
			s.error(w, r, http.StatusInternalServerError, s.available.reason)
			s.logger.Errorf("service is dead reason: %", s.available.reason)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// recoverPanic panic recovering middleware
func (s *server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := s.logger.WithFields(logrus.Fields{
					"remout_addr": r.RemoteAddr,
					"request_id":  r.Context().Value(ctxKeyRequestID),
					"method":      r.Method,
					"URI":         r.RequestURI,
				})

				logger.Errorf("ended hadling by panic with error: %s", err)

				s.error(w, r, http.StatusInternalServerError, ErrPanicHanding)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
