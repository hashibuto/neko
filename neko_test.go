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
