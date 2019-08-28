// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"net/http"
	"strings"
	"time"
)

// repoBase returns the first segment of the requested URL. It does so
// by splitting on `/`, and returning the first result.
func repoBase(url string) string {
	return strings.Split(url[1:], "/")[0]
}

// checkUpstream verifies that the package is available on the remote server
func (s *Server) checkUpstream(module string) (bool, int) {
	base := repoBase(module)
	upstream := s.repo.root + "/" + base
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Head will follow up to 10 redirects, so no need to worry about
	// it here.
	resp, err := client.Head(upstream)
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
func (s *Server) getRedirect(module string) string {
	base := repoBase(module)
	if s.redirect == s.repo.root {
		return s.repo.root + "/" + base
	} else {
		return s.redirect + module
	}
}
