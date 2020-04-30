package vanity

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneric(t *testing.T) {
	mock := mockServer(t)
	defer mock.Close()

	addr := mockAddr(mock)

	s, _ := NewServer("base", addr, "")
	s.client = mock.Client()

	handler := http.HandlerFunc(s.handleGeneric)

	checks := []struct {
		repo    string
		get     bool
		fakeGet bool
		code    int
	}{
		{"valid", true, false, 200},
		{"valid", true, true, 302},
		{"valid", false, false, 302},
		{"valid", false, true, 302},
		{"invalid", true, false, 404},
		{"invalid", true, true, 404},
		{"invalid", false, false, 404},
		{"invalid", false, true, 404},
		// Root redirect
		{"", true, false, 302},
		{"", true, true, 302},
		{"", false, false, 302},
		{"", false, true, 302},
	}

	for _, c := range checks {
		addr := "/" + c.repo
		if c.get {
			if c.fakeGet {
				addr += "?go-get=0"
			} else {
				addr += "?go-get=1"
			}
		}

		req, err := http.NewRequest("GET", addr, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if c.code != rr.Code {
			t.Errorf("Handler(%v) returned wrong status code; expected %v, got %v",
				addr, c.code, rr.Code)
		}
	}
}

func TestNewServer(t *testing.T) {
	checks := []struct {
		base string
		root string
		err  string
	}{
		{"", "", "Missing or invalid base value"},
		{"", "root", "Missing or invalid base value"},
		{"base", "", "Missing or invalid root value"},
	}

	for _, c := range checks {
		_, err := NewServer(c.base, c.err, "")

		if err != nil && err.Error() != c.err {
			t.Errorf("Mismatch in error message, expected '%v', got '%v'",
				c.err, err.Error())
		}
	}
}

func TestWebRoot(t *testing.T) {
	checks := []struct {
		dir string
		err error
	}{
		{"/non-existent", fmt.Errorf("stat /non-existent: no such file or directory")},
		{"/dev/null", fmt.Errorf("Web root /dev/null is not a directory")},
		{"/", nil},
	}

	s, _ := NewServer("base", "root", "")
	for _, c := range checks {
		err := s.WebRoot(c.dir)
		if err != c.err && err.Error() != c.err.Error() {
			t.Errorf("Mismatch in error response; expected %v, got %v", c.err, err)
		}
	}
}

func TestRootRedirect(t *testing.T) {
	s, _ := NewServer("base", "root", "")
	s.RootRedirect("https://github.com/nirenjan")
	s.Repo().SetProvider("github")
	s.QueryRemote(false)

	handler := http.HandlerFunc(s.handleGeneric)

	req, err := http.NewRequest("GET", "http://localhost:2369/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("Unexpected status code %v, expected %v", rr.Code, http.StatusFound)
	}

	if redir, ok := rr.Result().Header["Location"]; !ok || redir[0] != "https://github.com/nirenjan" {
		t.Errorf("Unexpected headers %#v", rr.Result().Header)
	}
}

func TestWellKnown(t *testing.T) {
	s, _ := NewServer("base", "root", "")
	s.WebRoot("/")

	handler := http.HandlerFunc(s.handleWellKnown)
	checks := []struct {
		path string
		code int
	}{
		{"/non-existent", http.StatusNotFound},
		{"/dev", http.StatusNotFound},
		{"/dev/null", http.StatusOK},
	}

	for _, c := range checks {
		req, err := http.NewRequest("GET", c.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if c.code != rr.Code {
			t.Errorf("Unexpected status code %v, expected %v", rr.Code, c.code)
		}
	}
}

func TestGetHandler(t *testing.T) {
	dummy := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
	}

	checks := []struct {
		method string
		code   int
	}{
		{"GET", http.StatusOK},
		{"POST", http.StatusMethodNotAllowed},
		{"PUT", http.StatusMethodNotAllowed},
		{"PATCH", http.StatusMethodNotAllowed},
		{"HEAD", http.StatusMethodNotAllowed},
		{"OPTIONS", http.StatusMethodNotAllowed},
	}

	for _, c := range checks {
		req, err := http.NewRequest(c.method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		getHandler(dummy)(rr, req)

		if c.code != rr.Code {
			t.Errorf("Unexpected status code %v, expected %v", rr.Code, c.code)
		}
	}
}
