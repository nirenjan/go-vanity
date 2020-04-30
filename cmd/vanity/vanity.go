package main // import nirenjan.org/vanity/cmd/vanity

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"nirenjan.org/vanity"
)

// Flags for server
var base, root, redirect, provider, vcs, rootRedirect, webRoot string
var listenTCP, listenUnix string
var noQueryRemote bool

func main() {
	flag.StringVar(&base, "base", "", "Base URL for vanity server (required)")
	flag.StringVar(&root, "root", "", "Root URL for VCS host (required)")
	flag.StringVar(&redirect, "redirect", "", "Redirect URL for browsers")
	flag.StringVar(&provider, "provider", "", "VCS Provider")
	flag.StringVar(&vcs, "vcs", "", "VCS type (git, subversion, etc.)")
	flag.StringVar(&rootRedirect, "root-redirect", "", "Redirect for requests to base URL")

	flag.StringVar(&webRoot, "web-root", "", "Directory containing the .well-known folder")
	flag.StringVar(&listenTCP, "listen-tcp", "", "Port to listen on for HTTP server")
	flag.StringVar(&listenUnix, "listen-unix", "", "Socket to listen on for HTTP server")
	flag.BoolVar(&noQueryRemote, "no-query-remote", false, "Don't query the remote server for repo presence")
	flag.Parse()

	stderr := log.New(os.Stderr, "vanity: ", 0)

	if base == "" {
		stderr.Fatal("Missing Base URL on command line")
	}

	if root == "" {
		stderr.Fatal("Missing Root URL on command line")
	}

	if listenTCP != "" && listenUnix != "" {
		stderr.Fatal("Conflicting arguments -listen-tcp and -listen-unix")
	}

	server, err := vanity.NewServer(base, root, redirect)
	if err != nil {
		stderr.Fatal(err)
	}

	if rootRedirect != "" {
		server.RootRedirect(rootRedirect)
	}

	if listenTCP != "" {
		l, err := net.Listen("tcp", listenTCP)
		if err != nil {
			stderr.Fatal(err)
		}
		server.Listen(l)
	}

	if listenUnix != "" {
		l, err := net.Listen("unix", listenUnix)
		if err != nil {
			stderr.Fatal(err)
		}
		server.Listen(l)
		defer os.Remove(listenUnix)
	}

	if webRoot != "" {
		if err := server.WebRoot(webRoot); err != nil {
			stderr.Fatal(err)
		}
	}

	if provider != "" {
		if err := server.Repo().SetProvider(provider); err != nil {
			stderr.Fatal(err)
		}
	}

	if vcs != "" {
		if err := server.Repo().SetType(vcs); err != nil {
			stderr.Fatal(err)
		}
	}

	server.QueryRemote(!noQueryRemote)

	log.Println("Starting vanity server")
	if listenTCP != "" {
		log.Println("Listening on", listenTCP)
	} else if listenUnix != "" {
		log.Println("Listening on", listenUnix)
	}
	log.Println("Base URL:", base)
	log.Println("Root URL:", root)
	if redirect != "" {
		log.Println("Redirect to:", redirect)
	}
	if rootRedirect != "" {
		log.Println("Redirect Root:", rootRedirect)
	}

	if provider != "" {
		log.Println("Provider:", provider)
	}

	if vcs != "" {
		log.Println("VCS Type:", vcs)
	}

	log.Println("Query Remote:", !noQueryRemote)
	if webRoot != "" {
		log.Println("Web root:", webRoot)
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

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
