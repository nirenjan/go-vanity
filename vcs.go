// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"fmt"
	"strings"
)

// This file manages the VCS structure

// VcsProvider is an enumeration used to identify known VCS providers
// such as GitHub, Bitbucket, etc.
type VcsProvider int

const (
	Unknown   = (VcsProvider)(iota) // Unknown provider, call SetTemplates directly
	GitHub                          // Github.com or Github Enterprise
	GitLab                          // Gitlab.com or Gitlab CE/EE
	Gitea                           // Gitea
	Gogs                            // Gogs
	Bitbucket                       // Bitbucket Cloud or Bitbucket Server
)

// SetRoot configures the root directory of the hosting provider where the
// package is hosted.
func (v *Vcs) SetRoot(r string) {
	v.root = r
}

// SetProvider configures the Vcs structure to use the corresponding provider
func (v *Vcs) SetProvider(provider VcsProvider) {
	switch provider {
	case Unknown:
		// Do nothing

	case GitHub, GitLab:
		v.vcsType = "git"
		v.dirFormat = "tree/master{/dir}"
		v.fileFormat = "blob/master{/dir}/{file}#L{line}"

	case Bitbucket, Gogs, Gitea:
		// Default vcsType for Bitbucket is git, since Bitbucket is
		// sunsetting the mercurial repositories.
		v.vcsType = "git"
		v.dirFormat = "src/master{/dir}"
		v.fileFormat = "src/master{/dir}/{file}#L{line}"

	default:
		panic(fmt.Sprintf("Unknown provider %v", v))
	}
}

// SetType sets the version control system type.
// It can be one of the following case-insensitive strings:
// Bazaar, Fossil, Git, Mercurial, Subversion
func (v *Vcs) SetType(t string) {
	switch strings.ToLower(t) {
	case "bazaar":
		v.vcsType = "bzr"

	case "fossil":
		v.vcsType = "fossil"

	case "git":
		v.vcsType = "git"

	case "mercurial":
		v.vcsType = "hg"

	case "subversion":
		v.vcsType = "svn"

	default:
		panic(fmt.Sprintf("Unknown VCS type %v", t))
	}
}

// SetTemplates sets the URL templates for the directory and file
// listings. These are used by godoc to map the identifiers back to
// the source listings.
func (v *Vcs) SetTemplates(dir, file string) {
	// Check the file template, if it is not empty, it should
	// contain at least one instance of {file}.
	if file != "" && !strings.Contains(file, "{file}") {
		panic(fmt.Sprintf("Invalid file template %v", file))
	}

	v.dirFormat = dir
	v.fileFormat = file
}
