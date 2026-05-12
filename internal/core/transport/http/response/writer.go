package core_http_response

import "net/http"

var statusCodeUninitialized = -1

type Writer struct {
	http.ResponseWriter
	statusCode int
}

func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: w,
		statusCode:     statusCodeUninitialized,
	}
}

func (rw *Writer) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *Writer) StatusCodeOrPanic() int {
	if rw.statusCode == statusCodeUninitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
