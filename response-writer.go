package neko

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

// NewResponseWriter returns a new response writer which wraps the response status code header writer, so that it
// is accessible to the application once written
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, false}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader == true {
		return
	}

	rw.wroteHeader = true
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) WroteHeader() bool {
	return rw.wroteHeader
}
