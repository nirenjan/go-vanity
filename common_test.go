// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRepoBase(t *testing.T) {
	checks := []struct {
		in  string
		out string
	}{
		{"semver/core", "semver"},
		{"semver", "semver"},
		{"/semver/core", "semver"},
		{"/semver", "semver"},
		{"", ""},
	}

	for _, c := range checks {
		res := repoBase(c.in)

		if c.out != res {
			t.Errorf("Mismatch in repoBase, expected %#v, got %#v", c.out, res)
		}
	}
}

func TestGetRedirect(t *testing.T) {
	checks := []struct {
		module   string
		root     string
		redirect string
		exp      string
	}{
		{"/semver/core", "github.com/nirenjan/", "godoc.org/nirenjan.org/", "godoc.org/nirenjan.org/semver/core"},
		{"/semver", "github.com/nirenjan/", "godoc.org/nirenjan.org/", "godoc.org/nirenjan.org/semver"},
		{"/semver/core", "github.com/nirenjan/", "", "github.com/nirenjan/semver"},
		{"/semver", "github.com/nirenjan/", "", "github.com/nirenjan/semver"},
	}

	for _, c := range checks {
		s, _ := NewServer("nirenjan.org", c.root, c.redirect)

		if res := s.getRedirect(c.module); c.exp != res {
			t.Errorf("Mismatch in Server.getRedirect, expected %#v, got %#v",
				c.exp, res)
		}
	}
}

func TestCheckUpstream(t *testing.T) {
	checks := []struct {
		query string
		ok    bool
		code  int
	}{
		{"/valid", true, http.StatusOK},
		{"/invalid", false, http.StatusNotFound},
		{"/timeout", false, http.StatusServiceUnavailable},
		{"/error", false, http.StatusInternalServerError},
		{"/loop", false, http.StatusServiceUnavailable},
	}

	mock := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/valid":
			w.Write([]byte("valid-repo"))

		case "/timeout":
			time.Sleep(time.Second / 8)
			w.Write([]byte("timeout data"))

		case "/error":
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		case "/loop":
			http.Redirect(w, r, "/loop", http.StatusFound)

		default:
			http.NotFound(w, r)
		}
	}))
	l, _ := net.Listen("tcp", "127.0.0.1:")
	mock.Listener = l
	mock.Start()
	defer mock.Close()

	s, _ := NewServer("base", "http://"+mock.Listener.Addr().String()+"/", "")
	s.client = mock.Client()
	s.client.Timeout = time.Second / 10

	for _, c := range checks {
		ok, code := s.checkUpstream(c.query)
		if ok != c.ok || code != c.code {
			t.Errorf("Mismatch in Server.checkUpstream(%v); expected (%v, %v), got (%v, %v)",
				c.query, c.ok, c.code, ok, code)
		}
	}

	// Try disabling queryRemote
	s.queryRemote = false
	for _, c := range checks {
		ok, code := s.checkUpstream(c.query)
		if !ok || code != http.StatusOK {
			t.Errorf("Mismatch in Server.checkUpstream(%v); queryRemote = false, got (%v, %v)",
				c.query, ok, code)
		}
	}
}
