// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

var shutdownToken string

func init() {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	shutdownToken = id.String()
	log.Printf("Shutdown token: %v\n", shutdownToken)
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

// getHandler restricts requests to use the GET handler
func getHandler(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.RawPath)
		if r.Method == http.MethodGet {
			h(w, r)
			return
		}

		// Respond with the error message
		methodNotAllowed(w, r)
	}
}

// unAuthorized returns a 401 error
func unAuthorized(w http.ResponseWriter, r *http.Request) {
	tplData := struct {
		Code    int
		Message string
	}{
		http.StatusUnauthorized,
		"Unauthorized",
	}

	var b strings.Builder

	errTemplate.Execute(&b, tplData)
	http.Error(w, b.String(), http.StatusUnauthorized)
}

// shutDown shuts down the given HTTP server
func shutDown(s *http.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.RawPath)
		if r.Method != http.MethodPost {
			methodNotAllowed(w, r)
			return
		}

		if auth, ok := r.Header["Shutdown-Token"]; !ok || auth[0] != shutdownToken {
			unAuthorized(w, r)
			return
		}

		w.Write([]byte("OK\n"))
		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
