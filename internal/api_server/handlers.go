package apiserver

import (
	"net/http"
)

func (s *server) handlePing() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "ok")
	})
}
