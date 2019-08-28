// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"testing"
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
