package neko

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashibuto/oof"
)

//go:embed VERSION
var VERSION string

var statusErr = &StatusErr{}

const routeKey = "nekoRouteKey"

const nekoBanner = " _ __   ___| | _____  \n| '_ \\ / _ \\ |/ / _ \\ \n| | | |  __/   < (_) |\n|_| |_|\\___|_|\\_\\___/"

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
	fmt.Printf("%s  %s\n\n", nekoBanner, VERSION)
	fmt.Printf("Listening @%s\n", n.Server.Addr)
	return n.Server.ListenAndServe()
}

// ServeTLS initiates a blocking call which serves connections over TLS until interrupted
func (n *Neko) ServeTLS(certFile string, keyFile string) error {
	fmt.Printf("%s  %s\n\n", nekoBanner, VERSION)
	fmt.Printf("Listening (TLS) @%s\n", n.Server.Addr)
	return n.Server.ListenAndServeTLS(certFile, keyFile)
}

// UnwrapStatusError unwraps err as a status error if it contains one, or returns nil
func (n *Neko) UnwrapStatusError(err error) *StatusErr {
	var sErr *StatusErr
	var oofErr *oof.OofError
	if errors.As(err, &sErr) {
		return sErr
	}

	if errors.As(err, &oofErr) {
		if errors.As(oofErr.OrigError, &sErr) {
			return sErr
		}
	}
	return nil
}

func (n *Neko) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := n.router.FindMatch(r.URL.Path)
	if route == nil {
		w.WriteHeader(404)
		return
	}

	err := route.ServeHTTP(NewResponseWriter(w), r.WithContext(context.WithValue(r.Context(), routeKey, route)))
	if err != nil {
		statusErr := n.UnwrapStatusError(err)
		if statusErr != nil {
			w.WriteHeader(statusErr.StatusCode)
		} else {
			w.WriteHeader(500)
		}
	}
}
