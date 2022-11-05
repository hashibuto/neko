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

// func TestRouterMatch(t *testing.T) {
// 	router := NewRouter()
// 	router.AddRoute("/v1/api/literature/{id}")
// }
