package main // import nirenjan.org/vanity/cmd/vanity

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nirenjan.org/vanity"
)

// Flags for server
var base, root, redirect, provider, vcs, root_redirect, web_root string
var port uint

func validateCommandLine() {
	if base == "" {
		log.Fatal("Missing Base URL on command line")
	}

	if root == "" {
		log.Fatal("Missing Root URL on command line")
	}

	if port > 65535 {
		log.Fatal("Invalid port value", port, "(must be < 65536)")
	}

	// Make sure the web-root directory exist, if it is given
	if web_root != "" {
		stat, err := os.Stat(web_root)
		if err != nil {
			log.Fatal(err)
		}

		if !stat.Mode().IsDir() {
			log.Fatal("Web root", web_root, "is not a directory")
		}
	}
}

func main() {
	const DefaultPort uint = 2369
	flag.StringVar(&base, "base", "", "Base URL for vanity server (required)")
	flag.StringVar(&root, "root", "", "Root URL for VCS host (required)")
	flag.StringVar(&redirect, "redirect", "", "Redirect URL for browsers")
	flag.StringVar(&provider, "provider", "", "VCS Provider")
	flag.StringVar(&vcs, "vcs", "", "VCS type (git, subversion, etc.)")
	flag.StringVar(&root_redirect, "root-redirect", "", "Redirect for requests to base URL")

	flag.StringVar(&web_root, "web-root", "", "Directory containing the .well-known folder")
	flag.UintVar(&port, "port", DefaultPort, "Port to listen for HTTP server")
	flag.Parse()

	validateCommandLine()

	server, err := vanity.NewServer(base, root, redirect)
	if err != nil {
		log.Fatal(err)
	}

	if root_redirect != "" {
		server.RootRedirect(root_redirect)
	}

	if port != DefaultPort {
		server.Listen(uint16(port))
	}

	if web_root != "" {
		server.WebRoot(web_root)
	}

	if provider != "" {
		server.Repo().SetProvider(provider)
	}

	if vcs != "" {
		server.Repo().SetType(vcs)
	}

	log.Println("Starting vanity server on port", port)
	log.Println("Base URL:", base)
	log.Println("Root URL:", root)
	if root_redirect != "" {
		log.Println("Redirect URL:", redirect)
	}

	if provider != "" {
		log.Println("Provider:", provider)
	}

	if vcs != "" {
		log.Println("VCS Type:", vcs)
	}

	if web_root != "" {
		log.Println("Web root:", web_root)
	}

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
