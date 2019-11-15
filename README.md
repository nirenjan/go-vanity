Vanity - a tiny handler for Golang package vanity URLs
------------------------------------------------------

[![GoDoc](https://godoc.org/nirenjan.org/vanity?status.svg)](https://godoc.org/nirenjan.org/vanity)

Vanity is a small server that handles vanity URLs for Golang packages. It
redirects requests of the form `/xyz` to `https://github.com/<user>/xyz`, as
well as handling the corresponding `go-get` parameter in the request URL.

It works by sending a `GET` request to the upstream server, and if found, it
will send a 302 redirect back. If the upstream sends back an error message, it
will return a 404 error to the client.

In addition, Vanity also supports the creation of `go-source` meta tags. This
allows tools like godoc.org to link to the sources.

# Library

Vanity is available as a library to be integrated into an application.

`import "nirenjan.org/vanity"`

# CLI tool

Vanity is also available as a command-line tool that leverages the library.

`go get nirenjan.org/vanity/cmd/vanity`

## CLI arguments

```
Usage of vanity:
  -base string
        Base URL for vanity server (required)
  -listen-tcp string
        Port to listen on for HTTP server
  -listen-unix string
        Socket to listen on for HTTP server
  -no-query-remote
        Don't query the remote server for repo presence
  -provider string
        VCS Provider
  -redirect string
        Redirect URL for browsers
  -root string
        Root URL for VCS host (required)
  -root-redirect string
        Redirect for requests to base URL
  -vcs string
        VCS type (git, subversion, etc.)
  -web-root string
        Directory containing the .well-known folder (defaults to $PWD)
```

### Example

```
vanity \
    -base nirenjan.org \
    -root https://github.com/nirenjan/go- \
    -redirect https://godoc.org/nirenjan.org/ \
    -root-redirect https://github.com/nirenjan/?tab=repositories&language=go \
    -provider github
```

