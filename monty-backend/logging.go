package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ResponseWriterWrapper struct {
	w      *http.ResponseWriter
	body   *bytes.Buffer
	status *int
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriterWrapper(w)

		defer func() {
			ctime := time.Now()
			year := ctime.Year()
			month := ctime.Month()
			day := ctime.Day()
			hour := ctime.Hour()
			minute := ctime.Minute()

			tstr := fmt.Sprintf("%d-%d-%d %d:%d", year, month, day, hour, minute)

			log.Printf("[%s] %s %s %s - %d", tstr, r.RemoteAddr, r.Method, r.URL, *rw.status)
		}()

		next.ServeHTTP(rw, r)
	})
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	var buf bytes.Buffer
	var status int = 200

	return &ResponseWriterWrapper{w: &w, body: &buf, status: &status}
}

func (rw *ResponseWriterWrapper) Header() http.Header {
	return (*rw.w).Header()
}

func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return (*rw.w).Write(b)
}

func (rw *ResponseWriterWrapper) WriteHeader(status int) {
	*rw.status = status
	(*rw.w).WriteHeader(status)
}

func (rw *ResponseWriterWrapper) String() string {
	var buf bytes.Buffer

	buf.WriteString("Response:")
	buf.WriteString("Headers:")

	for k, v := range (*rw.w).Header() {
		buf.WriteString(fmt.Sprintf("%s: %s", k, v))
	}

	buf.WriteString(fmt.Sprintf("Status: %d", *rw.status))

	buf.WriteString("Body:")
	buf.WriteString(rw.body.String())

	return buf.String()
}
