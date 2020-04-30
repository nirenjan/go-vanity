// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"html/template"
	"net"
	"net/http"
)

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

	// provider is the VCS provider platform, e.g. Github.
	provider string
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

	// queryRemote is a flag that enables checking if the remote repository
	// exists, and setting it to false will force the server to always return a
	// 200 or 302 code, even if the repository doesn't exist on the remote.
	queryRemote bool

	// listener is the port/socket on which to listen to. The default
	// is tcp:2369
	listener net.Listener

	// listenerInit is a flag that indicates if the listener has been
	// initialized. It is used when the user has not created a custom
	// listener, and the server must fallback to the default listener.
	listenerInit bool

	// root is the location to redirect the request to the root node "/".
	// This defaults to repo.root, but it may be overridden by RootRedirect
	rootRedirect string

	// template is used by the server to save the template pointer.
	// This is used by handleGeneric to return the formatted data.
	template *template.Template

	// httpServer is a reference to the HTTP server, it is used during
	// initial bringup and final shutdown.
	httpServer *http.Server

	// client is a reference to the HTTP client used for querying the
	// upstream server. A default client is created when the server is
	// initialized, but it can be swapped with a separate client for
	// test purposes.
	client *http.Client
}
