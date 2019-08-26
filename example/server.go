// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"nirenjan.org/vanity"
)

func main() {
	server := vanity.Server{
		BaseURL: "nirenjan.org",
		Repo:    "https://git.nirenjan.com/go/",
	}

	server.Serve()
}
