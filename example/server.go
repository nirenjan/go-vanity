// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"nirenjan.org/vanity"
)

func main() {
	server, _ := vanity.NewServer("nirenjan.org", "https://github.com/nirenjan/go-", "")
	server.Repo().SetProvider(vanity.GitHub)
	server.Serve()
}
