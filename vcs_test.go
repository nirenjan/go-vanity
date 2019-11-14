// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"testing"
)

func TestSetRoot(t *testing.T) {
	checks := []string{
		"https://github.com/nirenjan/",
		"https://git.nirenjan.com/go/",
		"https://git.nirenjan.com/n",
	}

	for _, v := range checks {
		var vcs Vcs

		vcs.SetRoot(v)

		if vcs.root != v {
			t.Errorf("Mismatch in Vcs.root, expected %#v, got %#v", v, vcs.root)
		}
	}
}

func TestSetProvider(t *testing.T) {
	checks := []struct {
		provider   string
		vcsType    string
		dirFormat  string
		fileFormat string
	}{
		{"GitHub", "git", "tree/master{/dir}", "blob/master{/dir}/{file}#L{line}"},
		{"GitLab", "git", "tree/master{/dir}", "blob/master{/dir}/{file}#L{line}"},
		{"Gitea", "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
		{"Gogs", "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
		{"Bitbucket", "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
	}

	for _, p := range checks {
		var vcs Vcs
		err := vcs.SetProvider(p.provider)

		if err != nil {
			t.Errorf("Expected nil, got error %v", err)
		}

		if p.vcsType != vcs.vcsType {
			t.Errorf("Mismatch in Vcs.vcsType, expected %#v, got %#v", p.vcsType, vcs.vcsType)
		}

		if p.dirFormat != vcs.dirFormat {
			t.Errorf("Mismatch in Vcs.dirFormat, expected %#v, got %#v", p.dirFormat, vcs.dirFormat)
		}

		if p.fileFormat != vcs.fileFormat {
			t.Errorf("Mismatch in Vcs.fileFormat, expected %#v, got %#v", p.fileFormat, vcs.fileFormat)
		}
	}
}

func TestSetProviderInvalid(t *testing.T) {
	var vcs Vcs

	err := vcs.SetProvider("unknown")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestSetType(t *testing.T) {
	checks := []struct {
		name  string
		ident string
	}{
		{"bazaar", "bzr"},
		{"Bazaar", "bzr"},
		{"baZaar", "bzr"},
		{"fossil", "fossil"},
		{"git", "git"},
		{"mercurial", "hg"},
		{"subversion", "svn"},
	}

	for _, c := range checks {
		var v Vcs

		err := v.SetType(c.name)

		if err != nil {
			t.Errorf("Expected nil, got error %v", err)
		}

		if c.ident != v.vcsType {
			t.Errorf("Mismatch in Vcs.vcsType, expected %#v, got %#v", c.ident, v.vcsType)
		}
	}
}

func TestSetTypeInvalid(t *testing.T) {
	var vcs Vcs

	err := vcs.SetType("cvs")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestSetTemplates(t *testing.T) {
	checks := []struct {
		dir  string
		file string
	}{
		{"", ""},
		{"", "{file}"},
		{"{/dir}", "{file}"},
		{"{/dir}", "{file}#L{line}"},
		{"{dir}", ""},
	}

	for _, c := range checks {
		var v Vcs

		err := v.SetTemplates(c.dir, c.file)

		if err != nil {
			t.Errorf("Expected nil, got error %v", err)
		}

		if c.dir != v.dirFormat {
			t.Errorf("Mismatch in Vcs.dirFormat, expected %#v, got %#v", c.dir, v.dirFormat)
		}

		if c.file != v.fileFormat {
			t.Errorf("Mismatch in Vcs.fileFormat, expected %#v, got %#v", c.file, v.fileFormat)
		}
	}
}

func TestSetTemplatesInvalid(t *testing.T) {
	var vcs Vcs

	err := vcs.SetTemplates("", "c")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
