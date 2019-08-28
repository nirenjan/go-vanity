Vanity - a tiny handler for Golang package vanity URLs
------------------------------------------------------

Vanity is a small server that handles vanity URLs for Golang packages. It
redirects requests of the form `/xyz` to `https://git.nirenjan.com/go/xyz`, as
well as handling the corresponding `go-get` parameter in the request URL.

It works by sending a `GET` request to the upstream server, and if found, it
will send a 302 redirect back. If the upstream sends back an error message, it
will return a 404 error to the client.

In addition, Vanity also supports the creation of `go-source` meta tags. This
allows tools like godoc.org to link to the sources.

