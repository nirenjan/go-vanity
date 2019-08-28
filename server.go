// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NewServer takes in a base URL to bind the vanity name, the root at which all
// packages are hosted, and an optional redirect value to redirect web
// browsers. If redirect is empty, then the value of root is used as the
// redirect value. This returns a *Server which is used to configure and serve
// the repository.
func NewServer(base, root, redirect string) (*Server, error) {
	// Trim any trailing slashes, this will simplify the template
	// handling later
	base = strings.TrimSuffix(base, "/")
	root = strings.TrimSuffix(root, "/")
	redirect = strings.TrimSuffix(redirect, "/")

	if base == "" {
		return nil, fmt.Errorf("Missing or invalid base value")
	}
	if root == "" {
		return nil, fmt.Errorf("Missing or invalid root value")
	}

	// Create a new Server object
	s := new(Server)

	// Copy the values to the server
	s.base = base
	s.repo.SetRoot(root)

	if redirect == "" {
		redirect = root
	}
	s.redirect = redirect

	// Set defaults for the new server object
	s.repo.SetType("git")
	s.webRoot = "./"
	s.listenPort = 2369

	// Set the template
	s.buildTemplate()

	return s, nil
}

// Repo returns the pointer to the repo Vcs object so that it can
// be configured by the application
func (s *Server) Repo() *Vcs {
	return &s.repo
}

// WebRoot changes the web root for serving the `/.well-known/` folder
func (s *Server) WebRoot(wr string) {
	s.webRoot = wr
}

// Listen changes the listening port for the *Server
func (s *Server) Listen(port uint16) {
	s.listenPort = port
}

// Serve serves the given vanity name as configured by the *Server object
func (s *Server) Serve() {
	http.HandleFunc("/.well-known/", s.handleWellKnown)
	http.HandleFunc("/", s.handleGeneric)

	port := fmt.Sprintf(":%v", s.listenPort)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleWellKnown handles the "/.well-known/" directory and serves files
// from it.
func (s *Server) handleWellKnown(w http.ResponseWriter, r *http.Request) {
	file := filepath.Join(s.webRoot, r.URL.Path)

	// Check if the file exists
	info, err := os.Stat(file)
	if os.IsNotExist(err) || info.IsDir() {
		// Not found, or is a directory
		http.NotFound(w, r)
		return
	}

	// We have the file, serve it using http
	http.ServeFile(w, r, file)
}

// handleGeneric handles the regular endpoint
func (s *Server) handleGeneric(w http.ResponseWriter, r *http.Request) {
	// Get the path to the requested image
	module := r.URL.EscapedPath()

	// Make sure that the upstream exists
	exists, _ := s.checkUpstream(module)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Check if we got go-get=1 in the query
	redirect := func(r *http.Request) bool {
		get, ok := r.URL.Query()["go-get"]
		fmt.Printf("%#v %#v\n", ok, get)
		if !ok {
			// go-get was not in the query
			// Redirect to the redirect URL
			return true
		}

		// Search all the values for a matching one
		for _, v := range get {
			if v == "1" {
				return false
			}
		}

		return true
	}
	if redirect(r) {
		http.Redirect(w, r, s.getRedirect(module), http.StatusFound)
		return
	}

	s.serveMeta(w, module)
}
