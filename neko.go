package neko

import (
	"fmt"
	"net/http"
)

const VERSION = "v0.1.0"

type Neko struct {
	Server *http.Server

	// These two are corelated by index
	routePaths    []string
	routeHandlers []Handler
}

// NewServer returns a new Neko instance
func NewServer(server *http.Server) (*Neko, error) {
	if server.Handler != nil {
		return nil, fmt.Errorf("Neko specifies its own handler stack, server.Handler must be nil")
	}

	return &Neko{
		Server: server,
	}, nil
}

func (n *Neko) AddRoute(routePath string, handler Handler) {
	n.routePaths = append(n.routePaths, routePath)
	n.routeHandlers = append(n.routeHandlers, handler)
}

// Serve initiates a blocking call which serves connections until interrupted
func (n *Neko) Serve() error {
	fmt.Printf(" _ __   ___| | _____  \n| '_ \\ / _ \\ |/ / _ \\ \n| | | |  __/   < (_) |\n|_| |_|\\___|_|\\_\\___/  %s\n\n", VERSION)
	fmt.Printf("Listening @%s\n", n.Server.Addr)
	return n.Server.ListenAndServe()
}
