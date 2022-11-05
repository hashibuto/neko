package neko

import (
	"net/http"
	"testing"
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

func TestServe(t *testing.T) {
	s, err := NewServer(&http.Server{
		Addr: "localhost:8888",
	})
	if err != nil {
		t.Error(err)
		return
	}

	s.Serve()
}
