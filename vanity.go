// Package vanity creates a server to respond to Go package Vanity URLs. It
// converts the request URL to a VCS URL using a defined algorithm, checks if
// the target exists, and returns the redirect using the `go-import` meta tags
package vanity

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Server is a configuration structure to adjust the attributes of the vanity
// URL responder
type Server struct {
	// BaseURL is the Base URL to which the vanity name is bound.  E.g., for
	// the vanity name `rsc.io/quote/v3`, the BaseURL is `rsc.io`. This is a
	// required parameter
	BaseURL string

	// Repo is the location of the repository where the packages are hosted.
	// This is prepended to the package name.  E.g., the vanity name
	// `rsc.io/quote/v3` is hosted at `https://github.com/rsc/quote`. This is a
	// required parameter
	Repo string

	// Format is the version control system identifier, e.g., `git`. If this
	// is empty, it defaults to `git`.
	Format string

	// Redirect is the location to redirect to when a user directly navigates
	// to the vanity URL. E.g., navigating to `https://rsc.io/quote/v3`
	// redirects to `https://godoc.org/rsc.io/quote/v3`. If this is empty,
	// it defaults to the contents of `GitRepo`
	Redirect string

	// WebRoot is the location where to serve the contents of `.well-known`
	// folder. If this is empty, it will default to the current working
	// directory. This is used to conform to RFC 8615.
	WebRoot string

	// Listen is the default port on which to listen to. This will only listen
	// on IPv4 localhost, eg. ":8080"
	Listen uint16
}

// setDefaults sets the default values for the Server structure
func (s *Server) setDefaults() {
	// Default for Format
	if s.Format == "" {
		s.Format = "git"
	}

	// Default for Redirect
	if s.Redirect == "" {
		s.Redirect = s.Repo
	}

	// Default for WebRoot
	if s.WebRoot == "" {
		s.WebRoot = "./"
	}

	// Default for Listen port
	if s.Listen == 0 {
		s.Listen = 2369
	}

	// Truncate any trailing slashes from BaseURL, Repo and Redirect
	// This avoid duplicate slashes in the generated output
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
	s.Repo = strings.TrimSuffix(s.Repo, "/")
	s.Redirect = strings.TrimSuffix(s.Redirect, "/")
}

// checkConfig verifies that the required arguments are set
func (s Server) checkConfig() error {
	if s.BaseURL == "" {
		return fmt.Errorf("Missing BaseURL")
	}

	if s.Repo == "" {
		return fmt.Errorf("Missing Repo")
	}

	return nil
}

// repoBase returns the first segment of the requested URL. It does so
// by splitting on `/`, and returning the first result.
func repoBase(url string) string {
	return strings.Split(url, "/")[0]
}

// getTemplate returns a compiled HTML template to use in the server
func getTemplate() *template.Template {
	const tpl = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="go-import" content="{{.PackageBase}}/{{.RepoBase}} {{.Format}} {{.RepoRoot}}/{{.RepoBase}}">
	<meta http-equiv="refresh" content="0;url={{.Redirect}}{{- if eq .Redirect .RepoRoot -}}/{{.RepoBase}}{{- else -}}{{.Request}}{{- end -}}">
</head>
</html>`

	return template.Must(template.New("vanity").Parse(tpl))
}

// checkDestination verifies that the remote server is available
func checkDestination(dest string) (bool, int) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Head will follow up to 10 redirects, so no need to worry about
	// it here.
	resp, err := client.Head(dest)
	if err != nil {
		return false, http.StatusServiceUnavailable
	}

	// Close the body, avoid leaking resources
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, resp.StatusCode
}

// getRedirect gets the URL to redirect to
// If s.Redirect and s.Repo are the same, we cannot use the full request
// and must use the base only.
func (s Server) getRedirect(module string) string {
	base := repoBase(module[1:])
	if s.Redirect == s.Repo {
		return s.Repo + "/" + base
	} else {
		return s.Redirect + module
	}
}

func (s Server) Serve() {
	// Make sure required attributes are set
	if err := s.checkConfig(); err != nil {
		log.Fatal(err.Error())
	}

	// Set defaults for any unset attributes
	s.setDefaults()

	// Get the compiled template for the request. This will avoid having
	// to parse the template for every request.
	tpl := getTemplate()

	// Handle the /.well-known endpoint
	http.HandleFunc("/.well-known/", func(w http.ResponseWriter, r *http.Request) {
		file := filepath.Join(s.WebRoot, r.URL.Path)

		// Check if the file exists
		info, err := os.Stat(file)
		if os.IsNotExist(err) || info.IsDir() {
			// Not found, or is a directory
			http.NotFound(w, r)
			return
		}

		// We have the file, serve it using http
		http.ServeFile(w, r, file)
	})

	// Handle the regular endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the path to the requested image
		module := r.URL.EscapedPath()
		base := repoBase(module[1:])
		upstream := s.Repo + "/" + base

		// Make sure that the upstream exists
		exists, _ := checkDestination(upstream)
		if !exists {
			http.NotFound(w, r)
			return
		}

		// Check if we got go-get=1 in the query
		get, ok := r.URL.Query()["go-get"]
		if !ok || len(get[0]) < 1 || get[0] != "1" {
			// go-get=1 was not in the query
			// Redirect to the redirect URL
			http.Redirect(w, r, s.getRedirect(module), http.StatusFound)
			return
		}

		// Build the template data
		tplData := struct {
			PackageBase string
			RepoBase    string
			Format      string
			RepoRoot    string
			Redirect    string
			Request     string
		}{
			PackageBase: s.BaseURL,
			RepoBase:    base,
			Request:     module,
			RepoRoot:    s.Repo,
			Redirect:    s.Redirect,
			Format:      s.Format,
		}

		// Execute the template and write to w
		tpl.Execute(w, tplData)
	})

	port := fmt.Sprintf(":%v", s.Listen)
	log.Fatal(http.ListenAndServe(port, nil))
}
