// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	if strings.HasSuffix(root, "/") {
		// Trim the trailing multiple slashes and add a single /
		root = strings.TrimSuffix(root, "/") + "/"
	}
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
	s.rootRedirect = root

	if redirect == "" {
		redirect = root
	}
	s.redirect = redirect

	// Set defaults for the new server object
	s.repo.SetType("git")
	s.webRoot = "./"

	s.queryRemote = true
	s.client = new(http.Client)
	s.client.Timeout = time.Second * 5

	// Set the template
	s.buildTemplate()

	return s, nil
}

// Repo returns the pointer to the repo Vcs object so that it can
// be configured by the application
func (s *Server) Repo() *Vcs {
	return &s.repo
}

// RootRedirect changes the redirect for the `/` endpoint
func (s *Server) RootRedirect(rr string) {
	s.rootRedirect = rr
}

// WebRoot changes the web root for serving the `/.well-known/` folder
func (s *Server) WebRoot(wr string) error {
	stat, err := os.Stat(wr)
	if err != nil {
		return err
	}

	if !stat.Mode().IsDir() {
		return fmt.Errorf("Web root %v is not a directory", wr)
	}

	s.webRoot = wr
	return nil
}

// Listen changes the listening port/socket for the *Server
func (s *Server) Listen(l net.Listener) {
	if l != s.listener {
		if s.listenerInit {
			s.listener.Close()
		}
		s.listener = l
		s.listenerInit = true
	}
}

// QueryRemote controls whether the server should query the remote for
// existence of the requested repository. By default, this is true, causing
// the server to return 404 if the remote doesn't exist. However, this can
// be disabled so that the server always assumes that the remote repo exists.
func (s *Server) QueryRemote(query bool) {
	s.queryRemote = query
}

// Serve serves the given vanity name as configured by the *Server object
func (s *Server) Serve() error {
	m := http.NewServeMux()
	s.httpServer = &http.Server{Handler: m}

	m.HandleFunc("/.well-known/", getHandler(s.handleWellKnown))
	m.HandleFunc("/robots.txt", getHandler(handleRobots))
	m.HandleFunc("/robots.txt/", getHandler(http.NotFound))
	m.HandleFunc("/", getHandler(s.handleGeneric))

	if !s.listenerInit {
		var err error
		s.listener, err = net.Listen("tcp", "127.0.0.1:2369")
		if err != nil {
			return err
		}
	}

	if err := s.httpServer.Serve(s.listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Printf("Finished")

	return nil
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

	// If the module is the root node, redirect to the root redirect
	if module == "/" {
		http.Redirect(w, r, s.rootRedirect, http.StatusFound)
		return
	}

	// Make sure that the upstream exists
	exists, _ := s.checkUpstream(module)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Check if we got go-get=1 in the query
	redirect := func(r *http.Request) bool {
		get, ok := r.URL.Query()["go-get"]

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
