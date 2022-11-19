package neko

import (
	"net/http"
	"testing"

	"github.com/hashibuto/oof"
)

type DummyHandler struct{}

func (h *DummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func TestSpecifiedHandler(t *testing.T) {
	_, err := NewServer(&http.Server{
		Handler: &DummyHandler{},
	})
	if err == nil {
		t.Error("Expected an error to be returned when supplying a handler to the server object")
	}
}

func TestStatusErrorUnwrap(t *testing.T) {
	e := NewStatusErrf(401, "Unauthorized")
	s, err := NewServer(&http.Server{
		Addr: "localhost:8888",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if s.UnwrapStatusError(e) == nil {
		t.Error("Should not be nil")
		return
	}
}

func TestStatusDoubleErrorUnwrap(t *testing.T) {
	e := oof.Trace(NewStatusErrf(401, "Unauthorized"))
	s, err := NewServer(&http.Server{
		Addr: "localhost:8888",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if s.UnwrapStatusError(e) == nil {
		t.Error("Should not be nil")
		return
	}
}

// func TestServe(t *testing.T) {
// 	s, err := NewServer(&http.Server{
// 		Addr: "localhost:8888",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	//s.Use(middleware.RequestLogger)
// 	s.Route("/v1/test/{id:int}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
// 		mapping := ParsePathTokens(r)
// 		fmt.Println(mapping)
// 		return NewStatusErrf(201, "hello")
// 	})

// 	s.Serve()
// }
