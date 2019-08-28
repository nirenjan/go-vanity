// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import "text/template"

// Vcs is a configuration structure to configure the version control system
// provider.
type Vcs struct {
	// root is the root directory of the VCS host where the packages are
	// hosted. E.g., the package `rsc.io/quote` is hosted under the root
	// `https://github.com/rsc/`
	root string

	// vcsType is the version control system identifier, e.g., `git`. The
	// default value is `git`.
	vcsType string

	// dirFormat is the URL template for the directory at the VCS host.
	// This is used by godoc to map the identifiers back to the source.
	dirFormat string

	// fileFormat is the URL template for the file at the VCS host.
	// This is used by godoc to map the identifiers back to the source.
	fileFormat string
}

// Server is a configuration structure to adjust the attributes of the vanity
// URL responder
type Server struct {
	// base is the base URL to which the vanity name is bound.  E.g., for
	// the vanity name `rsc.io/quote/v3`, the BaseURL is `rsc.io`.
	base string

	// repo is the VCS provider that hosts the package.
	// `rsc.io/quote/v3` is hosted at `https://github.com/rsc/quote`.
	repo Vcs

	// redirect is the location to redirect to when a user directly navigates
	// to the vanity URL. E.g., navigating to `https://rsc.io/quote/v3`
	// redirects to `https://godoc.org/rsc.io/quote/v3`. The value of redirect
	// in this example is `https://godoc.org/rsc.io`. If this is empty, it
	// defaults to the contents of `repo.root`.
	redirect string

	// webRoot is the location where to serve the contents of `.well-known`
	// folder. If this is empty, it will default to the current working
	// directory. This is used to conform to RFC 8615.
	webRoot string

	// listenPort is the default port on which to listen to. This will only
	// listen on IPv4 localhost, eg. ":8080". The default is 2369
	listenPort uint16

	// template is used by the server to save the template pointer.
	// This is used by handleGeneric to return the formatted data.
	template *template.Template
}