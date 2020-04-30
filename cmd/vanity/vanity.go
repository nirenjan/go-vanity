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

	logger := log.New(os.Stderr, "vanity: ", 0)
	validateArgs(logger)

	server := spawnServer(logger)
	configureServer(logger, server)

	log.SetFlags(0)
	log.Println("Starting vanity server")
	if listenTCP != "" {
		log.Println("Listening on", listenTCP)
	} else if listenUnix != "" {
		log.Println("Listening on", listenUnix)
	}
	log.Print(server)
	log.SetFlags(log.LstdFlags)

	// Handle os.Interrupt
	go func() {
		ch := make(chan os.Signal)

		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ch:
			if listenUnix != "" {
				os.Remove(listenUnix)
			}
			server.ShutDown()
		}
	}()

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

func validateArgs(logger *log.Logger) {
	if base == "" {
		logger.Fatal("Missing Base URL on command line")
	}

	if root == "" {
		logger.Fatal("Missing Root URL on command line")
	}

	if listenTCP != "" && listenUnix != "" {
		logger.Fatal("Conflicting arguments -listen-tcp and -listen-unix")
	}

}

func spawnServer(logger *log.Logger) *vanity.Server {
	server, err := vanity.NewServer(base, root, redirect)
	if err != nil {
		logger.Fatal(err)
	}

	if rootRedirect != "" {
		server.RootRedirect(rootRedirect)
	}

	if listenTCP != "" {
		l, err := net.Listen("tcp", listenTCP)
		if err != nil {
			logger.Fatal(err)
		}
		server.Listen(l)
	}

	if listenUnix != "" {
		l, err := net.Listen("unix", listenUnix)
		if err != nil {
			logger.Fatal(err)
		}
		server.Listen(l)
	}

	return server
}

func configureServer(logger *log.Logger, server *vanity.Server) {
	if webRoot != "" {
		if err := server.WebRoot(webRoot); err != nil {
			logger.Fatal(err)
		}
	}

	if provider != "" {
		if err := server.Repo().SetProvider(provider); err != nil {
			logger.Fatal(err)
		}
	}

	if vcs != "" {
		if err := server.Repo().SetType(vcs); err != nil {
			logger.Fatal(err)
		}
	}

	server.QueryRemote(!noQueryRemote)
}
