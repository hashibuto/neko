package neko

import (
	"net/http"
	"testing"
)

var DummyHandlerFunc = func(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestParseRouteTokenDoesntSpan(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Parser should have detected bad token")
		}
	}()
	router := NewRouter()
	router.AddRoute("/v1/api/d{id}", false)
}

func TestParseRouteTokenMalformed(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Parser should have detected malformed token")
		}
	}()
	router := NewRouter()
	router.AddRoute("/v1/api/{id-one}", false)

}

func TestParseRouteGoodToken(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Error(err)
		}
	}()
	router := NewRouter()
	router.AddRoute("/v1/api/{id}", false)
}

func TestParseRouteGoodTokenWithType(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Error(err)
		}
	}()
	router := NewRouter()
	router.AddRoute("/v1/api/{id:int}", false)
}

func TestRouterMatch(t *testing.T) {
	router := NewRouter()
	router.AddRoute("/v1/api/literature/{id}", false).HandlerFunc(DummyHandlerFunc, "GET", "POST")
	router.AddRoute("/v1/api/animals/dog/{id}", false).HandlerFunc(DummyHandlerFunc, "GET", "POST")
	router.AddRoute("/v1/api/animals/cat/{id}", false).HandlerFunc(DummyHandlerFunc)
	router.AddRoute("/v1/api/people/{job}/type/{type_id}", false).HandlerFunc(DummyHandlerFunc)
	router.AddRoute("/v1/api/people/{job}/group/{id}", false).HandlerFunc(DummyHandlerFunc)
	router.AddRoute("/static/", true).HandlerFunc(DummyHandlerFunc)

	route := router.FindMatch("/v1/api/animals/cat/1223-3389-4345-3445")
	if route == nil {
		t.Error("Expected a route match")
		return
	}

	if route.path != "/v1/api/animals/cat/{id}" {
		t.Errorf("Got incorrect match on %s", route.path)
		return
	}

	route = router.FindMatch("/v1/api/animals/cats/1223-3389-4345-3445")
	if route != nil {
		t.Error("Expected a non-match")
		return
	}

	route = router.FindMatch("/v1/api/people/1223-3389-4345-3445")
	if route != nil {
		t.Error("Expected a non-match")
		return
	}

	route = router.FindMatch("/static/domino.txt")
	if route == nil {
		t.Error("Expected a match")
		return
	}
}
