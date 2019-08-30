// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"nirenjan.org/vanity"
)

func main() {
	server, _ := vanity.NewServer("nirenjan.org", "https://github.com/nirenjan/go-", "")
	server.Repo().SetProvider(vanity.GitHub)
	server.RootRedirect("https://github.com/nirenjan?utf8=%E2%9C%93&tab=repositories&q=&type=&language=go")
	server.Serve()
}
