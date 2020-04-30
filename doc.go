// Copyright 2019 Nirenjan Krishnan. All rights reserved.

/*
Package vanity creates a server to respond to Go package Vanity URLs. It
converts the request URL to a VCS URL using a defined algorithm, checks if
the target exists, and returns the redirect using the `go-import` meta
tags, along with optional `go-source` meta tags.

The vanity server serves *insecure* HTTP only, and is intended to be used
behind a web server which proxies the requests to the vanity server. The
fronting web server is responsible for terminating HTTPS connections.
*/
package vanity // import "nirenjan.org/vanity"
