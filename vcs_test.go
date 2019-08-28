// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"testing"
)

func expectPanic(t *testing.T, f string) {
	if r := recover(); r == nil {
		t.Errorf("Function %#v did not panic", f)
	}
}

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
		provider   VcsProvider
		vcsType    string
		dirFormat  string
		fileFormat string
	}{
		{GitHub, "git", "tree/master{/dir}", "blob/master{/dir}/{file}#L{line}"},
		{Unknown, "", "", ""},
		{GitLab, "git", "tree/master{/dir}", "blob/master{/dir}/{file}#L{line}"},
		{Gitea, "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
		{Gogs, "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
		{Bitbucket, "git", "src/master{/dir}", "src/master{/dir}/{file}#L{line}"},
	}

	for _, p := range checks {
		var vcs Vcs
		vcs.SetProvider(p.provider)

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

func TestSetProviderPanic(t *testing.T) {
	var vcs Vcs

	defer expectPanic(t, "SetProvider")
	vcs.SetProvider(VcsProvider(0xFFFF))
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

		v.SetType(c.name)

		if c.ident != v.vcsType {
			t.Errorf("Mismatch in Vcs.vcsType, expected %#v, got %#v", c.ident, v.vcsType)
		}
	}
}

func TestSetTypePanic(t *testing.T) {
	var vcs Vcs

	defer expectPanic(t, "SetType")
	vcs.SetType("cvs")
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

		v.SetTemplates(c.dir, c.file)

		if c.dir != v.dirFormat {
			t.Errorf("Mismatch in Vcs.dirFormat, expected %#v, got %#v", c.dir, v.dirFormat)
		}

		if c.file != v.fileFormat {
			t.Errorf("Mismatch in Vcs.fileFormat, expected %#v, got %#v", c.file, v.fileFormat)
		}
	}
}

func TestSetTemplatesPanic(t *testing.T) {
	var vcs Vcs

	defer expectPanic(t, "SetTemplates")
	vcs.SetTemplates("", "c")
}
