package apiserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WrightHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader((statusCode))
}
