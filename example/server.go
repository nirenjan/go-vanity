// Vanity is a redirector for custom Golang vanity URLs
package main

import (
	"os"
	"os/signal"
	"syscall"

	"nirenjan.org/vanity"
)

func main() {
	server, _ := vanity.NewServer("nirenjan.org", "https://github.com/nirenjan/go-", "")
	server.Repo().SetProvider("GitHub")
	server.RootRedirect("https://github.com/nirenjan?utf8=%E2%9C%93&tab=repositories&q=&type=&language=go")

	// Handle os.Interrupt
	go func() {
		ch := make(chan os.Signal)

		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ch:
			server.ShutDown()
		}
	}()

	server.Serve()
}
