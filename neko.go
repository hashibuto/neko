package neko

import (
	"errors"
	"fmt"
	"net/http"
)

const VERSION = "v0.1.0"

var statusErr = &StatusErr{}

type Middleware func(next Handler) Handler

type Neko struct {
	Server *http.Server

	router *Router
}

// NewServer returns a new Neko instance
func NewServer(server *http.Server) (*Neko, error) {
	if server.Handler != nil {
		return nil, fmt.Errorf("Neko specifies its own handler stack, server.Handler must be nil")
	}

	n := &Neko{
		Server: server,
		router: NewRouter(),
	}
	server.Handler = n

	return n, nil
}

// Use applies the target middleware to all routes added after invocation of the command.  Routes added prior
// to invoking "Use" will not invoke the middleware.
func (n *Neko) Use(mw Middleware) {
	n.router.AddMiddleware(mw)
}

// Route adds a complete match route to the server, route matching is independent of the order in which the
// route was added
func (n *Neko) Route(routePath string) *Route {
	return n.router.AddRoute(routePath, false)
}

// RoutePrefix adds a prefix match route to the server, route matching is independent of the order in which the
// route was added
func (n *Neko) RoutePrefix(routePath string) *Route {
	return n.router.AddRoute(routePath, true)
}

// Serve initiates a blocking call which serves connections until interrupted
func (n *Neko) Serve() error {
	fmt.Printf(" _ __   ___| | _____  \n| '_ \\ / _ \\ |/ / _ \\ \n| | | |  __/   < (_) |\n|_| |_|\\___|_|\\_\\___/  %s\n\n", VERSION)
	fmt.Printf("Listening @%s\n", n.Server.Addr)
	return n.Server.ListenAndServe()
}

func (n *Neko) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := n.router.FindMatch(r.URL.Path)
	if route == nil {
		w.WriteHeader(404)
		return
	}

	respWriter := NewResponseWriter(w)
	err := route.ServeHTTP(respWriter, r)
	if err != nil {
		var sErr *StatusErr
		if errors.As(err, &sErr) {
			w.WriteHeader(sErr.StatusCode)
		} else {
			w.WriteHeader(500)
		}

	}
}
