// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var errTemplate *template.Template

func init() {
	errTemplate = template.Must(template.New("error").Parse(`
<html>
<head>
<title>{{ .Code }} {{ .Message }}</title>
</head>
<body bgcolor="white">
<center><h1>{{ .Code }} {{ .Message }}</h1></center>
<hr><center>nirenjan.org/vanity</center>
</body>
</html>
`))
}

// methodNotAllowed returns a not allowed response
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	tplData := struct {
		Code    int
		Message string
	}{
		http.StatusMethodNotAllowed,
		"Method Not Allowed",
	}

	var b strings.Builder

	errTemplate.Execute(&b, tplData)
	http.Error(w, b.String(), http.StatusMethodNotAllowed)
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
		log.Printf("%v %v -> %d %v", r.Method, r.URL.Path, lrw.statusCode, elapsed)
	}
}

// ShutDown shuts down the HTTP server and gracefully exits
func (s *Server) ShutDown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
