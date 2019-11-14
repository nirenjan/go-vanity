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
var base, root, redirect, provider, vcs, root_redirect, web_root string
var listen_tcp, listen_unix string
var noQueryRemote bool

func main() {
	flag.StringVar(&base, "base", "", "Base URL for vanity server (required)")
	flag.StringVar(&root, "root", "", "Root URL for VCS host (required)")
	flag.StringVar(&redirect, "redirect", "", "Redirect URL for browsers")
	flag.StringVar(&provider, "provider", "", "VCS Provider")
	flag.StringVar(&vcs, "vcs", "", "VCS type (git, subversion, etc.)")
	flag.StringVar(&root_redirect, "root-redirect", "", "Redirect for requests to base URL")

	flag.StringVar(&web_root, "web-root", "", "Directory containing the .well-known folder")
	flag.StringVar(&listen_tcp, "listen-tcp", "", "Port to listen on for HTTP server")
	flag.StringVar(&listen_unix, "listen-unix", "", "Socket to listen on for HTTP server")
	flag.BoolVar(&noQueryRemote, "no-query-remote", false, "Don't query the remote server for repo presence")
	flag.Parse()

	stderr := log.New(os.Stderr, "vanity: ", 0)

	if base == "" {
		stderr.Fatal("Missing Base URL on command line")
	}

	if root == "" {
		stderr.Fatal("Missing Root URL on command line")
	}

	if listen_tcp != "" && listen_unix != "" {
		stderr.Fatal("Conflicting arguments -listen-tcp and -listen-unix")
	}

	server, err := vanity.NewServer(base, root, redirect)
	if err != nil {
		stderr.Fatal(err)
	}

	if root_redirect != "" {
		server.RootRedirect(root_redirect)
	}

	if listen_tcp != "" {
		l, err := net.Listen("tcp", listen_tcp)
		if err != nil {
			stderr.Fatal(err)
		}
		server.Listen(l)
	}

	if listen_unix != "" {
		l, err := net.Listen("unix", listen_unix)
		if err != nil {
			stderr.Fatal(err)
		}
		server.Listen(l)
		defer os.Remove(listen_unix)
	}

	if web_root != "" {
		if err := server.WebRoot(web_root); err != nil {
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
	if listen_tcp != "" {
		log.Println("Listening on", listen_tcp)
	} else if listen_unix != "" {
		log.Println("Listening on", listen_unix)
	}
	log.Println("Base URL:", base)
	log.Println("Root URL:", root)
	if redirect != "" {
		log.Println("Redirect to:", redirect)
	}
	if root_redirect != "" {
		log.Println("Redirect Root:", root_redirect)
	}

	if provider != "" {
		log.Println("Provider:", provider)
	}

	if vcs != "" {
		log.Println("VCS Type:", vcs)
	}

	log.Println("Query Remote:", !noQueryRemote)
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

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
