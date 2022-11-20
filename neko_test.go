package neko

import (
	"fmt"
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
	if UnwrapStatusError(e) == nil {
		t.Error("Should not be nil")
		return
	}
}

func TestStatusDoubleErrorUnwrap(t *testing.T) {
	e := oof.Trace(oof.Trace(NewStatusErrf(401, "Unauthorized")))
	val := UnwrapStatusError(e)

	if val == nil {
		t.Error("Should not be nil")
		return
	}

	fmt.Println(val.StatusCode)
}
