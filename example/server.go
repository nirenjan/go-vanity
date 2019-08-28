// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"nirenjan.org/vanity"
)

func main() {
	server, _ := vanity.NewServer("nirenjan.org", "https://git.nirenjan.com/go", "")
	server.Serve()
}
