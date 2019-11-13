package vanity

import (
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
		repo     string
		get      bool
		fake_get bool
		code     int
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
			if c.fake_get {
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
