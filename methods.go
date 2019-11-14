// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"context"
	"log"
	"net/http"
	"time"
)

// methodNotAllowed returns a not allowed response
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	message := http.StatusText(code)

	http.Error(w, message, code)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// getHandler restricts requests to use the GET handler
func getHandler(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}

		if r.Method == http.MethodGet {
			h(lrw, r)
		} else {
			// Respond with the error message
			methodNotAllowed(lrw, r)
		}

		t := time.Now()
		elapsed := t.Sub(start)
		remote := r.Header.Get("X-Forwarded-For")
		if remote == "" {
			remote = r.RemoteAddr
		}
		log.Printf("%v \"%v %v %v\" %d %v", remote, r.Method, r.URL.Path, r.Proto, lrw.statusCode, elapsed)
	}
}

// ShutDown shuts down the HTTP server and gracefully exits
func (s *Server) ShutDown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
