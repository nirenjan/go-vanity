// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Configuration
const (
	// BaseURL is the Base URL to which the vanity name is bound.
	// E.g., for the vanity name nirenjan.org/vanity, the BaseURL
	// is `nirenjan.org`
	BaseURL = `nirenjan.org`

	// Requests of the form /abc are redirected to
	// fmt.Sprint(<RedirectFormat>, "abc")
	RedirectFormat = `https://git.nirenjan.com/go/%s`

	// Web root for placing .well-known files
	WebRoot = "./"
)

// checkDestination verifies that the remote server is available
func checkDestination(dest string) (bool, int) {
	client := &http.Client{
		Timeout: time.Second * 10,
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

// serveWellKnown will serve the file from .well-known
func serveWellKnown(w http.ResponseWriter, r *http.Request) {
	file := filepath.Join(WebRoot, r.URL.Path)

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

// sendRedirect will check for the presence of the go-get argument, and
// if it is equal to 1, then it will write the Go Import style template
// otherwise, it will redirect it to the upstream
func sendRedirect(w http.ResponseWriter, r *http.Request) {
	const metaTag string = `
<!DOCTYPE html>
<html>
<head>
    <meta name="go-import" content="%s git %s">
    <meta http-equiv="refresh" content="0;url=%s">
</head>
</html>
`
	// Get the path to the requested image
	module := r.URL.EscapedPath()[1:]
	base := BaseURL + r.URL.EscapedPath()
	upstream := fmt.Sprintf(RedirectFormat, module)

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
		// Redirect to upstream
		http.Redirect(w, r, upstream, http.StatusFound)
		return
	}

	// go-get=1 was in the query, return the meta tag
	fmt.Fprintf(w, metaTag, base, upstream, upstream)
}

func main() {
	http.HandleFunc("/.well-known/", serveWellKnown)
	http.HandleFunc("/", sendRedirect)
	log.Fatal(http.ListenAndServe(":2369", nil))
}
